package coaxialcontrolio

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func GetStatus(ctx context.Context, c dahuarpc.Conn, channel int) (Status, error) {
	res, err := dahuarpc.Send[struct {
		Status Status `json:"status"`
	}](ctx, c, dahuarpc.
		New("CoaxialControlIO.getStatus").
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

func GetCaps(ctx context.Context, c dahuarpc.Conn, channel int) (Caps, error) {
	res, err := dahuarpc.Send[struct {
		Caps Caps `json:"caps"`
	}](ctx, c, dahuarpc.
		New("CoaxialControlIO.getCaps").
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
	Type        Type        `json:"Type"`
	IO          IO          `json:"IO"`
	TriggerMode TriggerMode `json:"TriggerMode"`
}

func Control(ctx context.Context, c dahuarpc.Conn, channel int, controls ...ControlRequest) error {
	if len(controls) == 0 {
		return nil
	}

	_, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("CoaxialControlIO.control").
		Params(struct {
			Channel int              `json:"channel"`
			Info    []ControlRequest `json:"info"`
		}{
			Channel: channel,
			Info:    controls,
		}))

	return err
}

type Type int

const (
	TypeWhiteLight Type = 1
	TypeSpeaker    Type = 2
)

type IO int

const (
	On  IO = 1
	Off IO = 2
)

type TriggerMode int

const (
	TriggerModeLinked TriggerMode = 1
	TriggerModeManual TriggerMode = 2
)
