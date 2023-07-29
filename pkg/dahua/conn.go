package dahua

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Conn struct {
	client    *http.Client
	Camera    Camera
	State     State
	lastID    int
	Session   string
	Error     error
	LastLogin time.Time
}

func NewConn(client *http.Client, camera Camera) *Conn {
	return &Conn{
		client: client,
		Camera: camera,
		State:  StateLogout,
	}
}

func (c *Conn) RPC(ctx context.Context) (RequestBuilder, error) {
	if c.Session == "" {
		return RequestBuilder{}, ErrInvalidSession
	}

	return NewRequestBuilder(
		c.client,
		c.nextID(),
		c.Camera.rpc,
		c.Session,
	), nil
}

func (c *Conn) RPCLogin() RequestBuilder {
	return NewRequestBuilder(
		c.client,
		c.nextID(),
		c.Camera.rpcLogin,
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
	if c.State != StateLogout {
		panic("cannot set session when not logged out")
	}
	c.Session = session
}

func (c *Conn) Set(newState State, err ...error) {
	if c.State == StateLogout && newState == StateLogin {
		c.LastLogin = time.Now()
	} else if c.State == StateLogin && newState == StateLogin {
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

	c.State = newState
}

type Camera struct {
	ip       string
	rpc      string
	rpcLogin string
}

func NewCamera(ip string) Camera {
	return Camera{
		ip:       ip,
		rpc:      fmt.Sprintf("http://%s/RPC2", ip),
		rpcLogin: fmt.Sprintf("http://%s/RPC2_Login", ip),
	}
}