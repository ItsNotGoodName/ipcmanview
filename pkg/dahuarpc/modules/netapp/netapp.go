package netapp

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type NetInterface struct {
	ConnStatus      string `json:"ConnStatus"`
	Name            string `json:"Name"`
	NetCardName     string `json:"NetCardName"`
	PhysicalAddress string `json:"PhysicalAddress"`
	Speed           int    `json:"Speed"`
	SupportLongPoE  bool   `json:"SupportLongPoE"`
	Type            string `json:"Type"`
	Valid           bool   `json:"Valid"`
}

func GetNetInterfaces(ctx context.Context, c dahuarpc.Conn) ([]NetInterface, error) {
	res, err := dahuarpc.Send[struct {
		NetInterface []NetInterface `json:"netInterface"`
	}](ctx, c, dahuarpc.New("netApp.getNetInterfaces"))

	return res.Params.NetInterface, err
}
