package auth

import (
	"context"

	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
)

var _ dahua.GenRPC = (*Conn)(nil)

type Conn struct {
	conn     *dahua.Conn
	username string
	password string
}

func NewConn(conn *dahua.Conn, username, password string) Conn {
	return Conn{
		conn:     conn,
		username: username,
		password: password,
	}
}

func (c *Conn) Logout(ctx context.Context) {
	if c.conn.State() == dahua.StateLogin {
		Logout(ctx, c.conn)
	}
}

func (c *Conn) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	return RPC(ctx, c.conn, c.username, c.password)
}

// RPC generates a RPC request that is authenticated.
func RPC(ctx context.Context, conn *dahua.Conn, username, password string) (dahua.RequestBuilder, error) {
	switch conn.State() {
	case dahua.StateLogout:
		err := Login(ctx, conn, username, password)
		if err != nil {
			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateLogin:
		if err := KeepAlive(ctx, conn); err != nil {
			if conn.State() == dahua.StateLogout {
				return RPC(ctx, conn, username, password)
			}

			return dahua.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahua.StateError:
		return dahua.RequestBuilder{}, conn.Error
	}

	panic("unhandled connection state")
}
