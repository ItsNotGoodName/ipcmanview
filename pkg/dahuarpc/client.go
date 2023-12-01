package dahuarpc

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Client      *http.Client
	RPCURL      string
	RPCLoginURL string
}

type ClientData struct {
	State     State
	Session   string
	Error     error
	LastLogin time.Time
}

type ClientState struct {
	lastID    int
	State     State
	Session   string
	Error     error
	LastLogin time.Time
}

func NewClientState() ClientState {
	return ClientState{
		lastID:    0,
		State:     StateLogout,
		Session:   "",
		Error:     nil,
		LastLogin: time.Time{},
	}
}

func (s *ClientState) Data() ClientData {
	return ClientData{
		State:     s.State,
		Session:   s.Session,
		Error:     s.Error,
		LastLogin: s.LastLogin,
	}
}

func (s *ClientState) SetSession(session string) error {
	// Session can only be set if we are StateLogout
	switch s.State {
	case StateLogout:
	default:
		return fmt.Errorf("invalid previous state: %s", s.State)
	}

	s.Session = session

	return nil
}

func (s *ClientState) SetLogin() error {
	// Login can only be set if we are not StateError or StateClosed
	switch s.State {
	case StateError, StateClosed:
		return fmt.Errorf("invalid previous state: %s", s.State)
	default:
	}

	s.State = StateLogin
	s.LastLogin = time.Now()

	return nil
}

func (s *ClientState) SetError(err error) error {
	// Error can only be set if we are StateLogin
	switch s.State {
	case StateLogin:
	default:
		return fmt.Errorf("invalid previous state: %s", s.State)
	}

	s.State = StateError
	s.Error = err

	return nil
}

func (s *ClientState) SetLogout() error {
	// Logout can only be set if we are StateLogin or StateError
	switch s.State {
	case StateLogin, StateError:
	default:
		return fmt.Errorf("invalid previous state: %s", s.State)
	}

	s.State = StateLogout
	s.lastID = 0
	s.Session = ""
	s.Error = nil

	return nil
}

func (s *ClientState) SetClose() {
	s.State = StateClosed
}

func (s *ClientState) NextID() int {
	s.lastID++
	return s.lastID
}

func (s *ClientState) RawRPC(c Client) RequestBuilder {
	return NewRequestBuilder(
		c.Client,
		s.NextID(),
		c.RPCURL,
		s.Session,
	)
}

func (s *ClientState) RawRPCLogin(c Client) RequestBuilder {
	return NewRequestBuilder(
		c.Client,
		s.NextID(),
		c.RPCLoginURL,
		s.Session,
	)
}
