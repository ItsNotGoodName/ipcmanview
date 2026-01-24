package encode

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type SupportDynamicBitrate struct {
	Stream  string `json:"Stream"`
	Support bool   `json:"Support"`
}

type SupportAICoding struct {
	AICoding bool `json:"AICoding"`
}

type VideoEncodeDevices struct {
	MultiAudioEncode            int                     `json:"MultiAudioEncode"`
	RecordIndividualResolution  bool                    `json:"RecordIndividualResolution"`
	SupportAICoding             []SupportAICoding       `json:"SupportAICoding"`
	SupportDynamicBitrate       []SupportDynamicBitrate `json:"SupportDynamicBitrate"`
	SupportIndividualResolution bool                    `json:"SupportIndividualResolution"`
}

type Caps struct {
	MaxExtraStream     int                  `json:"MaxExtraStream"`
	VideoEncodeDevices []VideoEncodeDevices `json:"VideoEncodeDevices"`
}

func GetCaps(ctx context.Context, c dahuarpc.Conn, channel int) (Caps, error) {
	res, err := dahuarpc.Send[struct {
		Caps Caps `json:"caps"`
	}](ctx, c, dahuarpc.
		New("encode.getCaps").
		Params(struct {
			Channel int `json:"channel"`
		}{
			Channel: channel,
		}))

	return res.Params.Caps, err
}
