package dahua

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var _ GenRPC = (*Conn)(nil)
var _ GenRPCLogin = (*Conn)(nil)

type Conn struct {
	state       State
	client      *http.Client
	rpcURL      string
	rpcLoginURL string
	lastID      int
	Session     string
	Error       error
	LastLogin   time.Time
}

func NewConn(client *http.Client, ip string) *Conn {
	return &Conn{
		state:       StateLogout,
		client:      client,
		rpcURL:      fmt.Sprintf("http://%s/RPC2", ip),
		rpcLoginURL: fmt.Sprintf("http://%s/RPC2_Login", ip),
	}
}

func (c *Conn) RPC(ctx context.Context) (RequestBuilder, error) {
	if c.Session == "" {
		return RequestBuilder{}, ErrInvalidSession
	}

	return NewRequestBuilder(
		c.client,
		c.nextID(),
		c.rpcURL,
		c.Session,
	), nil
}

func (c *Conn) RPCLogin() RequestBuilder {
	return NewRequestBuilder(
		c.client,
		c.nextID(),
		c.rpcLoginURL,
		c.Session,
	)
}

func (c *Conn) nextID() int {
	c.lastID += 1
	return c.lastID
}

type State = int

const (
	StateLogout State = iota
	StateLogin
	StateError
)

func (c *Conn) UpdateSession(session string) {
	if c.state != StateLogout {
		panic("cannot set session when not logged out")
	}
	c.Session = session
}

func (c *Conn) State() State {
	return c.state
}

func (c *Conn) Set(newState State, err ...error) {
	if c.state == StateLogout && newState == StateLogin {
		c.LastLogin = time.Now()
	} else if c.state == StateLogin && newState == StateLogin {
		c.LastLogin = time.Now()
	} else {
		c.lastID = 0
		c.Session = ""
		c.Error = nil
		c.LastLogin = time.Time{}
	}

	if newState == StateError {
		if len(err) == 0 {
			panic("no error was supplied")
		}
		c.Error = err[0]
	}

	c.state = newState
}
