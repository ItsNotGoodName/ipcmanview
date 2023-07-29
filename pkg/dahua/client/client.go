package client

import "github.com/ItsNotGoodName/pkg/dahua"

type Client struct {
	Username string
	Password string
	Camera   dahua.Camera
	Conn     *dahua.Conn
}

func (c Client) RPC() dahua.RequestBuilder {
	return c.Camera.RPC(c.Conn)
}

func (c Client) RPCLogin() dahua.RequestBuilder {
	return c.Camera.RPCLogin(c.Conn)
}
