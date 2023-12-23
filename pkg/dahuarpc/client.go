package dahuarpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	ErrClientRequestFailed = errors.New("client request failed")
	ErrClientClosed        = errors.New("client closed")
)

const (
	StateLogout State = iota
	StateLogin
	StateError
	StateClosed
)

type State int

func (s State) String() string {
	switch s {
	case StateLogin:
		return "login"
	case StateLogout:
		return "logout"
	case StateError:
		return "error"
	case StateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// func (s State) Is(states ...State) bool {
// 	return slices.Contains(states, s)
// }

type client struct {
	Client      *http.Client
	Username    string
	Password    string
	RPCURL      string
	RPCLoginURL string

	sync.Mutex
	clientState
}

type clientState struct {
	lastID    int
	state     State
	session   string
	error     error
	lastLogin time.Time
}

func (c *client) DoRaw(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error) {
	var url string
	if rb.Login {
		url = c.RPCLoginURL
	} else {
		url = c.RPCURL
	}

	b, err := json.Marshal(rb.Request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrClientRequestFailed, err)
	}
	return resp.Body, nil
}

func (c *client) fill(rb RequestBuilder) RequestBuilder {
	c.lastID++
	rb = rb.ID(c.lastID).Session(c.session)
	if arrs, ok := rb.Request.Params.([]RequestBuilder); ok {
		for i := range arrs {
			c.lastID++
			arrs[i] = arrs[i].ID(c.lastID).Session(c.session)
		}
	}
	return rb
}

type clientLogin struct {
	client *client
}

func (c clientLogin) Do(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error) {
	return c.client.DoRaw(ctx, c.client.fill(rb))
}

func (c clientLogin) SetSession(session string) {
	c.client.session = session
}

func NewClient(httpClient *http.Client, httpAddress, username, password string) Client {
	return Client{
		client: &client{
			clientState: clientState{
				lastID:    0,
				state:     StateLogout,
				session:   "",
				error:     nil,
				lastLogin: time.Time{},
			},
			Client:      httpClient,
			Username:    username,
			Password:    password,
			RPCURL:      URL(httpAddress),
			RPCLoginURL: LoginURL(httpAddress),
		},
	}
}

type Client struct {
	client *client
}

func (c Client) clientLogin() clientLogin {
	return clientLogin{
		client: c.client,
	}
}

func (c Client) SessionRaw() string {
	c.client.Lock()
	session := c.client.session
	c.client.Unlock()
	return session
}

func (c Client) Session(ctx context.Context) (string, error) {
	c.client.Lock()
	err := c.readyConnection(ctx)
	if err != nil {
		c.client.Unlock()
		return "", err
	}
	session := c.client.session
	c.client.Unlock()
	return session, err
}

func (c Client) Do(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error) {
	c.client.Lock()
	err := c.readyConnection(ctx)
	if err != nil {
		c.client.Unlock()
		return nil, err
	}
	rb = c.client.fill(rb)
	c.client.Unlock()

	return c.client.DoRaw(ctx, rb)
}

func (c Client) Close(ctx context.Context) error {
	c.client.Lock()
	var err error
	if c.client.state == StateLogin {
		_, err = Logout(ctx, c.clientLogin())
	}
	c.client.state = StateClosed
	c.client.Unlock()
	return err
}

func (c Client) readyConnection(ctx context.Context) error {
	switch c.client.state {
	case StateLogout:
		if err := Login(ctx, c.clientLogin(), c.client.Username, c.client.Password); err != nil {
			var e *LoginError
			if errors.As(err, &e) {
				c.client.clientState = clientState{
					lastID:    0,
					state:     StateError,
					session:   "",
					error:     err,
					lastLogin: time.Time{},
				}
			}
			return err
		}
		c.client.clientState = clientState{
			lastID:    c.client.lastID,
			state:     StateLogin,
			session:   c.client.session,
			error:     nil,
			lastLogin: time.Now(),
		}

		return nil
	case StateLogin:
		if time.Now().Sub(c.client.lastLogin) > 60*time.Second {
			if _, err := KeepAlive(ctx, c.clientLogin()); err != nil {
				if !errors.Is(err, ErrClientRequestFailed) {
					c.client.clientState = clientState{
						lastID:    0,
						state:     StateLogout,
						session:   "",
						error:     nil,
						lastLogin: time.Time{},
					}

					return c.readyConnection(ctx)
				}
				return err
			}
			c.client.clientState = clientState{
				lastID:    c.client.lastID,
				state:     StateLogin,
				session:   c.client.session,
				error:     nil,
				lastLogin: time.Now(),
			}
		}

		return nil
	case StateError:
		return c.client.error
	case StateClosed:
		return ErrClientClosed
	}

	panic("unhandled connection state")
}

type ClientState struct {
	State     State
	Session   string
	Error     error
	LastLogin time.Time
}

func (c Client) State() ClientState {
	c.client.Lock()
	data := ClientState{
		State:     c.client.state,
		Session:   c.client.session,
		Error:     c.client.error,
		LastLogin: c.client.lastLogin,
	}
	c.client.Unlock()
	return data
}
