package dahuarpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrClientClosed = errors.New("client closed")
)

const clientKeepAliveTimeout = 60 * time.Second

type Config struct {
	ctx     context.Context
	onError func(err error)
}

type ConfigFunc func(c *Config)

func WithContext(ctx context.Context) ConfigFunc {
	return func(c *Config) {
		c.ctx = ctx
	}
}

func WithOnError(fn func(err error)) ConfigFunc {
	return func(c *Config) {
		c.onError = fn
	}
}

func clientLogError(address string) func(err error) {
	return func(err error) {
		slog.Error("", slog.String("address", address), slog.String("package", "dahuarpc"), slog.String("error", err.Error()))
	}
}

type ClientState struct {
	State     State
	LastID    int
	Error     error
	LastLogin time.Time
	LastRPC   time.Time
}

type clientState struct {
	State     State
	LastID    int
	Session   string
	Error     error
	LastLogin time.Time
	Ticker    *time.Ticker
}

func (s *clientState) NextID() int {
	s.LastID++
	return s.LastID
}

func (s *clientState) To(newState State, err ...error) {
	switch newState {
	case StateLogout:
		s.State = StateLogout
		s.LastID = 0
		s.Session = ""
		s.Error = nil
		s.Ticker.Stop()
	case StateLogin:
		s.State = StateLogin
		s.Error = nil
		s.LastLogin = time.Now()
		s.Ticker.Reset(clientKeepAliveTimeout)
	case StateError:
		s.State = StateError
		if len(err) > 0 {
			s.Error = err[0]
		} else {
			s.Error = errors.New("error not set")
		}
		s.Ticker.Stop()
	case StateClosed:
		s.State = StateClosed
		s.Ticker.Stop()
	default:
		panic(fmt.Sprintf("unknown state: %s", newState))
	}
}

func (s *clientState) SetSession(session string) {
	s.Session = session
}

type clientLogin struct {
	*clientState
	client      *http.Client
	rpcURL      string
	rpcLoginURL string
}

func (s clientLogin) Do(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error) {
	var urL string
	if rb.Login {
		urL = s.rpcLoginURL
	} else {
		urL = s.rpcURL
	}
	return DoRaw(ctx, rb.ID(s.NextID()).Session(s.Session), s.client, urL)
}

func NewClient(httpClient *http.Client, u *url.URL, username, password string, configFuncs ...ConfigFunc) Client {
	cfg := Config{
		ctx:     context.Background(),
		onError: clientLogError(u.String()),
	}

	for _, fn := range configFuncs {
		fn(&cfg)
	}

	c := Client{
		client:      httpClient,
		username:    username,
		password:    password,
		rpcURL:      URL(u),
		rpcLoginURL: LoginURL(u),
		onError:     cfg.onError,
		doneC:       make(chan struct{}),
		rpcCC:       make(chan chan clientRPC),
		stateC:      make(chan ClientState),
		sessionCC:   make(chan chan string),
		closeCC:     make(chan chan error),
	}

	go c.serve(cfg.ctx)

	return c
}

type Client struct {
	client      *http.Client
	username    string
	password    string
	rpcURL      string
	rpcLoginURL string
	onError     func(err error)

	doneC chan struct{}

	rpcCC     chan chan clientRPC
	stateC    chan ClientState
	sessionCC chan chan string
	closeCC   chan chan error
}

func (c Client) clientLogin(rpcState *clientState) clientLogin {
	return clientLogin{
		clientState: rpcState,
		client:      c.client,
		rpcURL:      c.rpcURL,
		rpcLoginURL: c.rpcLoginURL,
	}
}

func (c Client) checkError(err error) {
	if err != nil {
		c.onError(err)
	}
}

// serve can only be called once and returns when context is canceled or client is closed.
// It handles authenticating and keeping the connection alive.
// If authentication eror occurs, then it will enter and errored state.
func (c Client) serve(ctx context.Context) {
	defer close(c.doneC)

	state := clientState{
		LastID:    0,
		State:     StateLogout,
		Session:   "",
		Error:     nil,
		LastLogin: time.Time{},
		Ticker:    time.NewTicker(clientKeepAliveTimeout),
	}

	login := func() {
		err := Login(ctx, c.clientLogin(&state), c.username, c.password)
		if err != nil {
			var e *LoginError
			if errors.As(err, &e) {
				state.To(StateError, err)
			} else {
				c.checkError(err)
			}
		} else {
			state.To(StateLogin)
		}
	}

	logout := func() error {
		var closeErr error
		if state.State == StateLogin {
			_, err := Logout(ctx, c.clientLogin(&state))
			var respErr *ResponseError
			if errors.As(err, &respErr) && respErr.Type == ErrorTypeInvalidSession {
				closeErr = nil
			} else {
				closeErr = err
			}
		}
		state.To(StateClosed)
		return closeErr
	}

	lastRPC := time.Time{}

	for {
		select {
		case <-ctx.Done():
			// Logout
			c.checkError(logout())
			return
		case <-state.Ticker.C:
			if state.State == StateLogin {
				// KeepAlive
				if _, err := KeepAlive(ctx, c.clientLogin(&state)); err != nil {
					if !errors.Is(err, ErrRequestFailed) {
						state.To(StateLogout)
						// Login
						login()
					} else {
						c.checkError(err)
					}
				} else {
					state.To(StateLogin)
				}
			}
		case rpcC := <-c.rpcCC:
			if state.State == StateLogout {
				// Login
				login()
			}

			var reply clientRPC
			switch state.State {
			case StateLogin:
				reply = clientRPC{
					ID:      state.NextID(),
					Session: state.Session,
				}
			default:
				var err error
				if state.Error != nil {
					err = state.Error
				} else {
					err = fmt.Errorf("invalid state: %s", state.State)
				}
				reply = clientRPC{
					Error: err,
				}
			}
			rpcC <- reply

			lastRPC = time.Now()
		case sessionC := <-c.sessionCC:
			if state.State == StateLogout {
				// Login
				login()
			}

			sessionC <- state.Session
		case c.stateC <- ClientState{
			State:     state.State,
			LastID:    state.LastID,
			Error:     state.Error,
			LastLogin: state.LastLogin,
			LastRPC:   lastRPC,
		}:
		case closeC := <-c.closeCC:
			// Logout
			closeC <- logout()
			return
		}
	}
}

type clientRPC struct {
	ID      int
	Session string
	Error   error
}

func (c Client) Do(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error) {
	if rb.Login {
		return nil, fmt.Errorf("login request not supported")
	}

	rpcC := make(chan clientRPC, 1)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.doneC:
		return nil, ErrClientClosed
	case c.rpcCC <- rpcC:
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.doneC:
		return nil, ErrClientClosed
	case rpc := <-rpcC:
		if rpc.Error != nil {
			return nil, rpc.Error
		}
		return DoRaw(ctx, rb.ID(rpc.ID).Session(rpc.Session), c.client, c.rpcURL)
	}
}

func (c Client) close(ctx context.Context, wait bool) error {
	errC := make(chan error, 1)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.doneC:
		return nil
	case c.closeCC <- errC:
	}

	if wait {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.doneC:
			return nil
		case err := <-errC:
			return err
		}
	}
	return nil
}

// CloseNoWait closes the connection without waiting for it to fully close.
func (c Client) CloseNoWait(ctx context.Context) error {
	return c.close(ctx, false)
}

// Close fully closes the connection.
func (c Client) Close(ctx context.Context) error {
	return c.close(ctx, true)
}

func (c Client) State(ctx context.Context) ClientState {
	select {
	case <-ctx.Done():
		return ClientState{Error: ctx.Err()}
	case <-c.doneC:
		return ClientState{
			State: StateClosed,
			Error: ctx.Err(),
		}
	case state := <-c.stateC:
		return state
	}
}

func (c Client) Session(ctx context.Context) string {
	sessionC := make(chan string, 1)
	select {
	case <-ctx.Done():
		return ""
	case <-c.doneC:
		return ""
	case c.sessionCC <- sessionC:
	}

	select {
	case <-ctx.Done():
		return ""
	case <-c.doneC:
		return ""
	case session := <-sessionC:
		return session
	}
}
