package config

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetVideoAnalyseRules(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[VideoAnalyseRules], error) {
	return configmanager.GetConfig[VideoAnalyseRules](ctx, c, "VideoAnalyseRule", true)
}

type VideoAnalyseRules []VideoAnalyseRule

func (m VideoAnalyseRules) Validate() error {
	return nil
}

type VideoAnalyseRule struct {
	Enable      bool     `json:"Enable"`
	ID          int      `json:"Id"`
	Name        string   `json:"Name"`
	ObjectTypes []string `json:"ObjectTypes"`
	PtzPresetID int      `json:"PtzPresetId"`
	TrackEnable bool     `json:"TrackEnable"`
	Type        string   `json:"Type"`
}
