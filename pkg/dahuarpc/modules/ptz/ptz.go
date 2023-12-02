package ptz

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Conn interface {
	dahuarpc.Conn
	Session() string
}

type Params struct {
	Code string `json:"code"`
	Arg1 int    `json:"arg1"`
	Arg2 int    `json:"arg2"`
	Arg3 int    `json:"arg3"`
	Arg4 int    `json:"arg4"`
}

func Start(ctx context.Context, c *Client, channel int, params Params) error {
	instance, err := c.InstanceGet(ctx, channel)
	if err != nil {
		return err
	}

	req, err := c.RPCSEQ(ctx)
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, req.
		Method("ptz.start").
		Params(params).
		Object(instance.Result.Integer()))
	return err
}

func Stop(ctx context.Context, c *Client, channel int, params Params) error {
	instance, err := c.InstanceGet(ctx, channel)
	if err != nil {
		return err
	}

	req, err := c.RPCSEQ(ctx)
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, req.
		Method("ptz.stop").
		Params(params).
		Object(instance.Result.Integer()))
	return err
}

type Preset struct {
	Index int    `json:"Index"`
	Name  string `json:"Name"`
}

func GetPresets(ctx context.Context, c *Client, channel int, params Params) ([]Preset, error) {
	instance, err := c.InstanceGet(ctx, channel)
	if err != nil {
		return nil, err
	}

	req, err := c.RPCSEQ(ctx)
	if err != nil {
		return nil, err
	}

	res, err := dahuarpc.Send[struct {
		Presets []Preset `json:"presets"`
	}](ctx, req.
		Method("ptz.getPresets").
		Params(params).
		Object(instance.Result.Integer()))
	if err != nil {
		return nil, err
	}

	return res.Params.Presets, nil
}
