package dahua

import "fmt"

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

func (c Camera) RPC(conn *Conn) RequestBuilder {
	return NewRequestBuilder(
		conn.client,
		conn.NextID(),
		c.rpc,
		conn.Session,
	).RequireSession()
}

func (c Camera) RPCLogin(conn *Conn) RequestBuilder {
	return NewRequestBuilder(
		conn.client,
		conn.NextID(),
		c.rpcLogin,
		conn.Session,
	)
}
