package coaxialcontrolio

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func GetStatus(ctx context.Context, c dahuarpc.Client, channel int) (Status, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return Status{}, err
	}

	res, err := dahuarpc.Send[struct {
		Status Status `json:"status"`
	}](ctx, rpc.
		Method("CoaxialControlIO.getStatus").
		Params(struct {
			Channel int `json:"channel"`
		}{
			Channel: channel,
		}))

	return res.Params.Status, err
}

type Status struct {
	WhiteLight string `json:"WhiteLight"`
	Speaker    string `json:"Speaker"`
}

func GetCaps(ctx context.Context, c dahuarpc.Client, channel int) (Caps, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return Caps{}, err
	}

	res, err := dahuarpc.Send[struct {
		Caps Caps `json:"caps"`
	}](ctx, rpc.
		Method("CoaxialControlIO.getCaps").
		Params(struct {
			Channel int `json:"channel"`
		}{
			Channel: channel,
		}))

	return res.Params.Caps, err
}

type Caps struct {
	SupportControlFullcolorLight int `json:"SupportControlFullcolorLight"`
	SupportControlLight          int `json:"SupportControlLight"`
	SupportControlSpeaker        int `json:"SupportControlSpeaker"`
}

type ControlRequest struct {
	Type        int `json:"Type"`
	IO          int `json:"IO"`
	TriggerMode int `json:"TriggerMode"`
}

func Control(ctx context.Context, c dahuarpc.Client, channel int, controls ...ControlRequest) error {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, rpc.
		Method("CoaxialControlIO.control").
		Params(struct {
			Channel int              `json:"channel"`
			Info    []ControlRequest `json:"info"`
		}{
			Channel: channel,
			Info:    controls,
		}))

	return err
}
