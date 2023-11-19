package ptz

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Params struct {
	Code string `json:"code"`
	Arg1 int    `json:"arg1"`
	Arg2 int    `json:"arg2"`
	Arg3 int    `json:"arg3"`
	Arg4 int    `json:"arg4"`
}

func Start(ctx context.Context, c *Client, channel int, params Params) error {
	instance, err := c.Instance(ctx, channel)
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
	instance, err := c.Instance(ctx, channel)
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
