package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetRecord(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[Record], error) {
	return configmanager.GetConfig[Record](ctx, c, "Record", true)
}

type Record struct {
	Format        string `json:"Format"`
	HolidayEnable bool   `json:"HolidayEnable"`
	PreRecord     int    `json:"PreRecord"`
	Redundancy    bool   `json:"Redundancy"`
	SnapShot      bool   `json:"SnapShot"`
	Stream        int    `json:"Stream"`
}

func (c Record) Merge(js string) (string, error) {
	return "", fmt.Errorf("%w: Merge not implemented for 'Record'", errors.ErrUnsupported)
}

func (c Record) Validate() error {
	return nil
}
