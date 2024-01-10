package config

import (
	"context"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetVideoAnalyseRules(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[VideoAnalyseRules], error) {
	return configmanager.GetConfig[VideoAnalyseRules](ctx, c, "VideoAnalyseRule", true)
}

type VideoAnalyseRule struct {
	Class       string   `json:"Class"`
	Enable      bool     `json:"Enable"`
	ID          int      `json:"Id"`
	Name        string   `json:"Name"`
	ObjectTypes []string `json:"ObjectTypes"`
	PtzPresetID int      `json:"PtzPresetId"`
	TrackEnable bool     `json:"TrackEnable"`
	Type        string   `json:"Type"`
}

type VideoAnalyseRules []VideoAnalyseRule

func (c VideoAnalyseRules) Merge(js string) (string, error) {
	var err error
	for i, v := range c {
		prefix := strconv.Itoa(i) + "."
		js, err = configmanager.Merge(js, []configmanager.MergeOption{
			{Path: prefix + "Class", Value: v.Class},
			{Path: prefix + "Enable", Value: v.Enable},
			{Path: prefix + "Id", Value: v.ID},
			{Path: prefix + "Name", Value: v.Name},
			{Path: prefix + "ObjectTypes", Value: v.ObjectTypes},
			{Path: prefix + "PtzPresetId", Value: v.PtzPresetID},
			{Path: prefix + "TrackEnable", Value: v.TrackEnable},
			{Path: prefix + "Type", Value: v.Type},
		})
		if err != nil {
			return "", err
		}
	}
	return js, nil
}

func (c VideoAnalyseRules) Validate() error {
	return nil
}
