package intervideo

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func ManagerGetVersion(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Info struct {
			Onvif string `json:"Onvif"`
		} `json:"info"`
	}](ctx, c, dahuarpc.
		New("IntervideoManager.getVersion").
		Params(struct {
			Name string `json:"Name"`
		}{
			Name: "Onvif",
		}))
	if err != nil {
		return "", err
	}

	return res.Params.Info.Onvif, nil
}
