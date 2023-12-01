package ptz

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type clientData struct {
	sync.Mutex
	Instance *dahuarpc.Instance
	Session  string
	ID       int
}

func newClientData(session string) clientData {
	return clientData{
		Instance: dahuarpc.NewInstance("ptz.factory.instance"),
		Session:  session,
		ID:       0,
	}
}

type Client struct {
	conn Conn
	data clientData
}

func NewClient(conn Conn) *Client {
	return &Client{
		conn: conn,
		data: newClientData(""),
	}
}

func (c *Client) InstanceGet(ctx context.Context, channel int) (dahuarpc.Response[json.RawMessage], error) {
	c.data.Lock()
	res, err := c.data.Instance.Get(ctx, c.conn, strconv.Itoa(channel), nil)
	c.data.Unlock()

	return res, err
}

func (c *Client) RPCSEQ(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	c.data.Lock()
	session := c.conn.Session()
	if session != c.data.Session {
		c.data = newClientData(session)
	}

	seq := getSeq(session, c.data.ID)
	c.data.ID = getNextID(c.data.ID)

	rpc, err := c.conn.RPC(ctx)
	c.data.Unlock()

	if err != nil {
		return dahuarpc.RequestBuilder{}, err
	}

	return rpc.Seq(seq), nil
}
