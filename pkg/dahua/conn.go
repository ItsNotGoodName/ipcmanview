package dahua

import (
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

func NewConn(client *http.Client, ip string) *Conn {
	return &Conn{
		client: client,
		Camera: newCamera(ip),
		State:  StateLogout,
	}
}

func (c *Conn) RPC() (RequestBuilder, error) {
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

func (c *Conn) UpdateSession(session string) error {
	if c.State != StateLogout {
		return fmt.Errorf("cannot set session when not logged out")
	}
	c.Session = session
	return nil
}

func (c *Conn) SetError(err error) {
	c.Set(StateError)
	c.Error = err
}

func (c *Conn) Set(newState State) {
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

	c.State = newState
}

type Camera struct {
	ip       string
	rpc      string
	rpcLogin string
}

func newCamera(ip string) Camera {
	return Camera{
		ip:       ip,
		rpc:      fmt.Sprintf("http://%s/RPC2", ip),
		rpcLogin: fmt.Sprintf("http://%s/RPC2_Login", ip),
	}
}
