package config

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetGeneral(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[General], error) {
	return configmanager.GetConfig[General](ctx, c, "General", false)
}

type General struct {
	LocalNo           int    `json:"LocalNo"`
	LockLoginEnable   bool   `json:"LockLoginEnable"`
	LockLoginTimes    int    `json:"LockLoginTimes"`
	LoginFailLockTime int    `json:"LoginFailLockTime"`
	MachineName       string `json:"MachineName"`
	MaxOnlineTime     int    `json:"MaxOnlineTime"`
}

func (g General) Validate() error {
	return nil
}
