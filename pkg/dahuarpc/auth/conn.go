package auth

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

var _ dahuarpc.Client = (*Conn)(nil)

type Conn struct {
	conn     *dahuarpc.Conn
	username string
	password string
}

func NewConn(conn *dahuarpc.Conn, username, password string) Conn {
	return Conn{
		conn:     conn,
		username: username,
		password: password,
	}
}

func (c *Conn) Logout(ctx context.Context) {
	if c.conn.State() == dahuarpc.StateLogin {
		Logout(ctx, c.conn)
	}
}

func (c *Conn) RPC(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	return RPC(ctx, c.conn, c.username, c.password)
}

// RPC generates a RPC request that is authenticated.
func RPC(ctx context.Context, conn *dahuarpc.Conn, username, password string) (dahuarpc.RequestBuilder, error) {
	switch conn.State() {
	case dahuarpc.StateLogout:
		err := Login(ctx, conn, username, password)
		if err != nil {
			return dahuarpc.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahuarpc.StateLogin:
		if err := KeepAlive(ctx, conn); err != nil {
			if conn.State() == dahuarpc.StateLogout {
				return RPC(ctx, conn, username, password)
			}

			return dahuarpc.RequestBuilder{}, err
		}

		return conn.RPC(ctx)
	case dahuarpc.StateError:
		return dahuarpc.RequestBuilder{}, conn.Error
	}

	panic("unhandled connection state")
}
