package config

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetNTP(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[NTP], error) {
	return configmanager.GetConfig[NTP](ctx, c, "NTP", false)
}

type NTP struct {
	Address  string `json:"Address"`
	Enable   bool   `json:"Enable"`
	Port     int    `json:"Port"`
	TimeZone int    `json:"TimeZone"`
	// TimeZoneDesc is the description of the TimeZone.
	TimeZoneDesc string `json:"TimeZoneDesc"`
	UpdatePeriod int    `json:"UpdatePeriod"`
}

func (c NTP) Merge(js string) (string, error) {
	return configmanager.Merge(js, []configmanager.MergeValues{
		{Path: "Address", Value: c.Address},
		{Path: "Enable", Value: c.Enable},
		{Path: "Port", Value: c.Port},
		{Path: "TimeZone", Value: c.TimeZone},
		{Path: "TimeZoneDesc", Value: c.TimeZoneDesc},
		{Path: "UpdatePeriod", Value: c.UpdatePeriod},
	})
}

func (c NTP) Validate() error {
	return nil
}
