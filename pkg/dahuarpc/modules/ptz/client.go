package ptz

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type ClientRPC interface {
	dahuarpc.Client
	Session() string
}

type clientData struct {
	Instance    *dahuarpc.Instance
	LastSession string
	ID          int
}

func newClientData() clientData {
	return clientData{
		Instance:    dahuarpc.NewInstance("ptz.factory.instance"),
		LastSession: "",
		ID:          0,
	}
}

type Client struct {
	rpc ClientRPC

	dataMu sync.Mutex
	data   clientData
}

func NewClient(clientRPC ClientRPC) *Client {
	return &Client{
		rpc:    clientRPC,
		dataMu: sync.Mutex{},
		data:   newClientData(),
	}
}

func (c *Client) Instance(ctx context.Context, channel int) (dahuarpc.Response[json.RawMessage], error) {
	c.dataMu.Lock()
	res, err := c.data.Instance.Get(ctx, c.rpc, strconv.Itoa(channel), nil)
	c.dataMu.Unlock()
	return res, err
}

func (c *Client) RPCSEQ(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	c.dataMu.Lock()
	session := c.rpc.Session()
	if session != c.data.LastSession {
		c.data = newClientData()
	}
	c.data.LastSession = session

	seq := getSeq(session, c.data.ID)
	c.data.ID = getNextID(c.data.ID)

	rpc, err := c.rpc.RPC(ctx)
	c.dataMu.Unlock()
	if err != nil {
		return dahuarpc.RequestBuilder{}, err
	}

	return rpc.Seq(seq), nil
}
