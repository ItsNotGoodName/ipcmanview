package license

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type LicenseInfo struct {
	AbroadInfo    string `json:"AbroadInfo"`
	AllType       bool   `json:"AllType"`
	DigitChannel  int    `json:"DigitChannel"`
	EffectiveDays int    `json:"EffectiveDays"`
	EffectiveTime int    `json:"EffectiveTime"`
	LicenseID     int    `json:"LicenseID"`
	ProductType   string `json:"ProductType"`
	Status        int    `json:"Status"`
	Username      string `json:"Username"`
}

func GetLicenseInfo(ctx context.Context, c dahuarpc.Client) ([]LicenseInfo, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return nil, err
	}

	res, err := dahuarpc.Send[[]struct {
		Info LicenseInfo `json:"Info"`
	}](ctx, rpc.Method("License.getLicenseInfo"))
	if err != nil {
		return nil, err
	}

	params := make([]LicenseInfo, len(res.Params))
	for i := range res.Params {
		params[i] = res.Params[i].Info
	}

	return params, nil
}
