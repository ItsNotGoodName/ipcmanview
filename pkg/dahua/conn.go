package dahua

import (
	"net/http"
	"time"
)

type State = int

type Conn struct {
	State     State
	client    *http.Client
	lastID    int
	Session   string
	Error     error
	LastLogin time.Time
}

func NewConn(client *http.Client) *Conn {
	return &Conn{
		client: client,
		State:  StateLogout,
	}
}

func (c *Conn) NextID() int {
	c.lastID += 1
	return c.lastID
}

const (
	StateLogout State = iota
	StateLogin
	StateError
)

func (c *Conn) SetSession(session string) error {
	if c.State != StateLogout {
		panic("cannot set session on logout")
	}
	c.Session = session
	return nil
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
