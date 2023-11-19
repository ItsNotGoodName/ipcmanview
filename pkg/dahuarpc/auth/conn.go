package auth

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"golang.org/x/sync/singleflight"
)

var ErrConnClosed = errors.New("conn closed")

type Conn struct {
	*dahuarpc.Conn
	group    *singleflight.Group
	username string
	password string
}

func NewConn(conn *dahuarpc.Conn, username, password string) *Conn {
	return &Conn{
		Conn:     conn,
		group:    &singleflight.Group{},
		username: username,
		password: password,
	}
}

func (c *Conn) RPC(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	resC := c.group.DoChan("auth.Conn.readyConnection", func() (interface{}, error) { return nil, c.readyConnection(ctx) })
	select {
	case <-ctx.Done():
		return dahuarpc.RequestBuilder{}, ctx.Err()
	case res := <-resC:
		if res.Err != nil {
			return dahuarpc.RequestBuilder{}, res.Err
		}

		return c.Conn.RawRPC(ctx)
	}
}

func (c *Conn) readyConnection(ctx context.Context) error {
	switch c.Conn.Data().State {
	case dahuarpc.StateLogout:
		err := Login(ctx, c.Conn, c.username, c.password)
		if err != nil {
			return err
		}

		return nil
	case dahuarpc.StateLogin:
		if err := KeepAlive(ctx, c.Conn); err != nil {
			if c.Conn.Data().State == dahuarpc.StateLogout {
				return c.readyConnection(ctx)
			}

			return err
		}

		return nil
	case dahuarpc.StateError:
		return c.Conn.Data().Error
	case dahuarpc.StateClosed:
		return ErrConnClosed
	}

	panic("unhandled connection state")
}

func (c *Conn) Close(ctx context.Context) error {
	return Close(ctx, c.Conn)
}

func (c *Conn) Login(ctx context.Context) error {
	return Login(ctx, c.Conn, c.username, c.password)
}

func (c *Conn) KeepAlive(ctx context.Context) error {
	return KeepAlive(ctx, c.Conn)
}
