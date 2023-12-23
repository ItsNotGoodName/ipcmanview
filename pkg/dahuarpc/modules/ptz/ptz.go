package ptz

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Conn interface {
	dahuarpc.Conn
	SessionRaw() string
}

func newClient(session string) *client {
	return &client{
		Mutex:   sync.Mutex{},
		Cache:   dahuarpc.NewCache(),
		Session: session,
		ID:      0,
	}
}

type client struct {
	sync.Mutex
	Cache   dahuarpc.Cache
	Session string
	ID      int
}

func NewClient(conn Conn) Client {
	return Client{
		conn:   conn,
		client: newClient(conn.SessionRaw()),
	}
}

type Client struct {
	conn   Conn
	client *client
}

func (c Client) Instance(ctx context.Context, channel int) (dahuarpc.Response[json.RawMessage], error) {
	c.client.Lock()
	res, err := c.client.Cache.Send(ctx, c.conn, strconv.Itoa(channel), dahuarpc.New("ptz.factory.instance"))
	c.client.Unlock()

	return res, err
}

func (c *Client) Seq(rb dahuarpc.RequestBuilder) dahuarpc.RequestBuilder {
	c.client.Lock()
	session := c.conn.SessionRaw()
	if session != c.client.Session {
		c.client = newClient(session)
	}

	seq := nextSeq(session, c.client.ID)
	c.client.ID = nextID(c.client.ID)
	c.client.Unlock()

	return rb.Option("seq", seq)
}

type Params struct {
	Code string `json:"code"`
	Arg1 int    `json:"arg1"`
	Arg2 int    `json:"arg2"`
	Arg3 int    `json:"arg3"`
	Arg4 int    `json:"arg4"`
}

func Start(ctx context.Context, c Client, channel int, params Params) error {
	instance, err := c.Instance(ctx, channel)
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, c.conn, c.Seq(dahuarpc.
		New("ptz.start").
		Params(params).
		Object(instance.Result.Integer())))
	return err
}

func Stop(ctx context.Context, c Client, channel int, params Params) error {
	instance, err := c.Instance(ctx, channel)
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, c.conn, c.Seq(dahuarpc.
		New("ptz.stop").
		Params(params).
		Object(instance.Result.Integer())))
	return err
}

type Preset struct {
	Index int    `json:"Index"`
	Name  string `json:"Name"`
}

func GetPresets(ctx context.Context, c Client, channel int, params Params) ([]Preset, error) {
	instance, err := c.Instance(ctx, channel)
	if err != nil {
		return nil, err
	}

	res, err := dahuarpc.Send[struct {
		Presets []Preset `json:"presets"`
	}](ctx, c.conn, c.Seq(dahuarpc.
		New("ptz.getPresets").
		Params(params).
		Object(instance.Result.Integer())))
	if err != nil {
		return nil, err
	}

	return res.Params.Presets, nil
}
