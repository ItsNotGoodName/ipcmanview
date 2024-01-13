package ptz

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

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

func NewClient(conn dahuarpc.ConnSession) Client {
	return Client{
		conn:   conn,
		client: newClient(conn.Session(context.Background())),
	}
}

type Client struct {
	conn   dahuarpc.ConnSession
	client *client
}

func (c Client) Instance(ctx context.Context, channel int) (dahuarpc.Response[json.RawMessage], error) {
	c.client.Lock()
	res, err := c.client.Cache.Send(ctx, c.conn, strconv.Itoa(channel), dahuarpc.New("ptz.factory.instance"))
	c.client.Unlock()

	return res, err
}

func (c *Client) Seq(ctx context.Context, rb dahuarpc.RequestBuilder) dahuarpc.RequestBuilder {
	c.client.Lock()
	session := c.conn.Session(ctx)
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

	_, err = dahuarpc.Send[any](ctx, c.conn, c.Seq(ctx, dahuarpc.
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

	_, err = dahuarpc.Send[any](ctx, c.conn, c.Seq(ctx, dahuarpc.
		New("ptz.stop").
		Params(params).
		Object(instance.Result.Integer())))
	return err
}

type Preset struct {
	Index int    `json:"Index"`
	Name  string `json:"Name"`
}

func GetPresets(ctx context.Context, c Client, channel int) ([]Preset, error) {
	instance, err := c.Instance(ctx, channel)
	if err != nil {
		return nil, err
	}

	res, err := dahuarpc.Send[struct {
		Presets []Preset `json:"presets"`
	}](ctx, c.conn, c.Seq(ctx, dahuarpc.
		New("ptz.getPresets").
		Object(instance.Result.Integer())))
	if err != nil {
		return nil, err
	}

	return res.Params.Presets, nil
}

type Status struct {
	Postion       [3]float64 `json:"Postion"`
	Action        string     `json:"Action"`
	ActionID      int        `json:"ActionID"`
	MoveStatus    string     `json:"MoveStatus"`
	TaskName      string     `json:"TaskName"`
	PanTiltStatus string     `json:"PanTiltStatus"`
}

func GetStatus(ctx context.Context, c Client, channel int) (Status, error) {
	instance, err := c.Instance(ctx, channel)
	if err != nil {
		return Status{}, err
	}

	res, err := dahuarpc.Send[struct {
		Status Status `json:"status"`
	}](ctx, c.conn, dahuarpc.
		New("ptz.getStatus").
		Object(instance.Result.Integer()))
	if err != nil {
		return Status{}, err
	}

	return res.Params.Status, nil
}
