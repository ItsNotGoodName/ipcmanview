package dahuarpc

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"
)

type State int

const (
	StateLogout State = iota
	StateLogin
	StateError
	StateClosed
)

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

func (s State) Is(states ...State) bool {
	return slices.Contains(states, s)
}

type ConnData struct {
	lastID    int
	State     State
	Session   string
	Error     error
	LastLogin time.Time
}

func (c *ConnData) nextID() int {
	c.lastID += 1
	return c.lastID
}

func NewConn(client *http.Client, httpAddress string) *Conn {
	return &Conn{
		client:      client,
		rpcURL:      fmt.Sprintf("%s/RPC2", httpAddress),
		rpcLoginURL: fmt.Sprintf("%s/RPC2_Login", httpAddress),
		dataMu:      sync.Mutex{},
		data: ConnData{
			State:     StateLogout,
			lastID:    0,
			Session:   "",
			Error:     nil,
			LastLogin: time.Time{},
		},
	}
}

type Conn struct {
	client      *http.Client
	rpcURL      string
	rpcLoginURL string

	dataMu sync.Mutex
	data   ConnData
}

func (c *Conn) RawRPC(ctx context.Context) (RequestBuilder, error) {
	c.dataMu.Lock()
	if c.data.Session == "" {
		c.dataMu.Unlock()
		return RequestBuilder{}, ErrInvalidSession
	}

	rpc := NewRequestBuilder(
		c.client,
		c.data.nextID(),
		c.rpcURL,
		c.data.Session,
	)
	c.dataMu.Unlock()

	return rpc, nil
}

func (c *Conn) RawRPCLogin() RequestBuilder {
	c.dataMu.Lock()
	rpc := NewRequestBuilder(
		c.client,
		c.data.nextID(),
		c.rpcLoginURL,
		c.data.Session,
	)
	c.dataMu.Unlock()

	return rpc
}

func (c *Conn) Data() ConnData {
	c.dataMu.Lock()
	data := c.data
	c.dataMu.Unlock()
	return data
}

func (c *Conn) Session() string {
	c.dataMu.Lock()
	session := c.data.Session
	c.dataMu.Unlock()
	return session
}

func (c *Conn) UpdateSession(session string) error {
	c.dataMu.Lock()
	if !c.data.State.Is(StateLogout, StateClosed) {
		c.dataMu.Unlock()
		return fmt.Errorf("cannot set session when logout or closed")
	}

	c.data.Session = session
	c.dataMu.Unlock()

	return nil
}

func (c *Conn) UpdateState(state State, err ...error) {
	c.dataMu.Lock()
	switch state {
	case StateLogout, StateClosed:
		// (*) => (Logout, Closed)
		c.data = ConnData{}
	case StateLogin:
		// (*) => (Login)
		if c.data.State.Is(StateLogout, StateLogin) {
			// (Logout, Login) => (Login)
			c.data.LastLogin = time.Now()
		}
	case StateError:
		// (*) => (Error)
		if len(err) == 0 {
			c.data.Error = fmt.Errorf("unknown error")
		} else {
			c.data.Error = err[0]
		}
	}

	c.data.State = state
	c.dataMu.Unlock()
}
