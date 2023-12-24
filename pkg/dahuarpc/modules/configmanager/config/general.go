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

func (g General) Merge(js string) (string, error) {
	return configmanager.Merge(js, []configmanager.MergeOption{
		{Path: "LocalNo", Value: g.LocalNo},
		{Path: "LockLoginEnable", Value: g.LockLoginEnable},
		{Path: "LockLoginTimes", Value: g.LockLoginTimes},
		{Path: "LoginFailLockTime", Value: g.LoginFailLockTime},
		{Path: "MachineName", Value: g.MachineName},
		{Path: "MaxOnlineTime", Value: g.MaxOnlineTime},
	})
}

func (g General) Validate() error {
	return nil
}
