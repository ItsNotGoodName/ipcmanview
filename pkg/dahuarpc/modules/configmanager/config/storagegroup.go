package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetStorageGroup(ctx context.Context, c dahuarpc.Conn) (configmanager.ConfigArray[StorageGroup], error) {
	return configmanager.GetConfigArray[StorageGroup](ctx, c, "StorageGroup")
}

type StorageGroup struct {
	Channels []struct {
		MaxPictures int         `json:"MaxPictures"`
		Path        interface{} `json:"Path"`
	} `json:"Channels"`
	FileHoldTime    int    `json:"FileHoldTime"`
	Memo            string `json:"Memo"`
	Name            string `json:"Name"`
	OverWrite       bool   `json:"OverWrite"`
	PicturePathRule string `json:"PicturePathRule"`
	RecordPathRule  string `json:"RecordPathRule"`
}

func (c StorageGroup) Merge(js string) (string, error) {
	return "", fmt.Errorf("%w: Merge not implemented for 'StorageGroup'", errors.ErrUnsupported)
}

func (c StorageGroup) Validate() error {
	return nil
}
