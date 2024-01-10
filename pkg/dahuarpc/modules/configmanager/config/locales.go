package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetLocales(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[Locales], error) {
	return configmanager.GetConfig[Locales](ctx, c, "Locales", false)
}

type Locales struct {
	DSTEnable bool `json:"DSTEnable"`
	DSTEnd    struct {
		Day    int `json:"Day"`
		Hour   int `json:"Hour"`
		Minute int `json:"Minute"`
		Month  int `json:"Month"`
		Week   int `json:"Week"`
		Year   int `json:"Year"`
	} `json:"DSTEnd"`
	DSTStart struct {
		Day    int `json:"Day"`
		Hour   int `json:"Hour"`
		Minute int `json:"Minute"`
		Month  int `json:"Month"`
		Week   int `json:"Week"`
		Year   int `json:"Year"`
	} `json:"DSTStart"`
	TimeFormat string `json:"TimeFormat"`
}

func (c Locales) Merge(js string) (string, error) {
	return "", fmt.Errorf("%w: Merge not implemented for 'Locales'", errors.ErrUnsupported)
}

func (c Locales) Validate() error {
	return nil
}
