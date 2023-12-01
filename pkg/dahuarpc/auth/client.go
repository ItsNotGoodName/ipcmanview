package auth

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/global"
	"golang.org/x/sync/singleflight"
)

var ErrClientClosed = errors.New("client closed")

type loginClient struct {
	client dahuarpc.Client
	state  *dahuarpc.ClientState
}

func (c loginClient) UpdateSession(session string) {
	c.state.SetSession(session)
}

func (c loginClient) RawRPC() dahuarpc.RequestBuilder {
	return c.state.RawRPC(c.client)
}

func (c loginClient) RawRPCLogin() dahuarpc.RequestBuilder {
	return c.state.RawRPCLogin(c.client)
}

type clientState struct {
	sync.Mutex
	dahuarpc.ClientState
}

func NewClient(httpClient *http.Client, httpAddress, username, password string) Client {
	return Client{
		client: dahuarpc.Client{
			Client:      httpClient,
			RPCURL:      dahuarpc.RPCURL(httpAddress),
			RPCLoginURL: dahuarpc.RPCLoginURL(httpAddress),
		},
		state: &clientState{
			Mutex:       sync.Mutex{},
			ClientState: dahuarpc.NewClientState(),
		},
		group:    &singleflight.Group{},
		username: username,
		password: password,
	}
}

type Client struct {
	client   dahuarpc.Client
	state    *clientState
	group    *singleflight.Group
	username string
	password string
}

func (c Client) loginClient() loginClient {
	return loginClient{
		client: c.client,
		state:  &c.state.ClientState,
	}
}

// RPCSession returns a RequestBuilder that is authorized.
func (c Client) RPC(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	resC := c.group.DoChan("auth.Client.ready", func() (interface{}, error) { return nil, c.ready(ctx) })
	select {
	case <-ctx.Done():
		return dahuarpc.RequestBuilder{}, ctx.Err()
	case res := <-resC:
		if res.Err != nil {
			return dahuarpc.RequestBuilder{}, res.Err
		}

		c.state.Lock()
		rpc := c.state.RawRPC(c.client)
		c.state.Unlock()

		return rpc, nil
	}
}

// RPCSession returns a session that is authorized.
func (c Client) RPCSession(ctx context.Context) (string, error) {
	resC := c.group.DoChan("auth.Client.ready", func() (interface{}, error) { return nil, c.ready(ctx) })
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-resC:
		return c.Session(), nil
	}
}

// Session returns the current session.
func (c Client) Session() string {
	c.state.Lock()
	session := c.state.Session
	c.state.Unlock()

	return session
}

func (c Client) Close(ctx context.Context) error {
	c.state.Lock()
	_, err := global.Logout(ctx, c.loginClient())
	c.state.SetClose()
	c.state.Unlock()

	return err
}

func (c Client) Data() dahuarpc.ClientData {
	c.state.Lock()
	data := c.state.Data()
	c.state.Unlock()

	return data
}

func (c Client) ready(ctx context.Context) error {
	c.state.Lock()
	err := c.readyConnection(ctx)
	c.state.Unlock()
	return err
}

func (c Client) readyConnection(ctx context.Context) error {
	switch c.state.State {
	case dahuarpc.StateLogout:
		err := global.Login(ctx, c.loginClient(), c.username, c.password)
		if err != nil {
			var e *global.LoginError
			if errors.As(err, &e) {
				if err := c.state.SetError(err); err != nil {
					return err
				}
			}
			return err
		}

		if err := c.state.SetLogin(); err != nil {
			return err
		}

		return nil
	case dahuarpc.StateLogin:
		if time.Now().Sub(c.state.LastLogin) > 60*time.Second {
			_, err := global.KeepAlive(ctx, c.loginClient())
			if err != nil {
				if !errors.Is(err, dahuarpc.ErrRequestFailed) {
					if err := c.state.SetLogout(); err != nil {
						return err
					}

					return c.readyConnection(ctx)
				}

				return err
			}

			if err := c.state.SetLogin(); err != nil {
				return err
			}
		}

		return nil
	case dahuarpc.StateError:
		return c.state.Error
	case dahuarpc.StateClosed:
		return ErrClientClosed
	}

	panic("unhandled connection state")
}
