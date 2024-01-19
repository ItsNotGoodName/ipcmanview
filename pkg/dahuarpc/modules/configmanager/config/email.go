package config

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
)

func GetEmail(ctx context.Context, c dahuarpc.Conn) (configmanager.Config[Email], error) {
	return configmanager.GetConfig[Email](ctx, c, "Email", false)
}

type Email struct {
	Address        string        `json:"Address"`
	Anonymous      bool          `json:"Anonymous"`
	AttachEnable   bool          `json:"AttachEnable"`
	Authentication bool          `json:"Authentication"`
	CustomTitle    []interface{} `json:"CustomTitle"`
	Enable         bool          `json:"Enable"`
	HealthReport   struct {
		Enable   bool `json:"Enable"`
		Interval int  `json:"Interval"`
	} `json:"HealthReport"`
	OnlyAttachment bool     `json:"OnlyAttachment"`
	Password       string   `json:"Password"`
	Port           int      `json:"Port"`
	Receivers      []string `json:"Receivers"`
	SendAddress    string   `json:"SendAddress"`
	SendInterv     int      `json:"SendInterv"`
	SslEnable      bool     `json:"SslEnable"`
	Title          string   `json:"Title"`
	TLSEnable      bool     `json:"TlsEnable"`
	UserName       string   `json:"UserName"`
}

func (c Email) Merge(js string) (string, error) {
	return configmanager.Merge(js, []configmanager.MergeValues{
		{Path: "Address", Value: c.Address},
		{Path: "Anonymous", Value: c.Anonymous},
		{Path: "AttachEnable", Value: c.AttachEnable},
		{Path: "Authentication", Value: c.Authentication},
		{Path: "CustomTitle", Value: c.CustomTitle},
		{Path: "Enable", Value: c.Enable},
		{Path: "HealthReport.Enable", Value: c.HealthReport.Enable},
		{Path: "HealthReport.Interval", Value: c.HealthReport.Interval},
		{Path: "OnlyAttachment", Value: c.OnlyAttachment},
		{Path: "Password", Value: c.Password},
		{Path: "Port", Value: c.Port},
		{Path: "Receivers", Value: c.Receivers},
		{Path: "SendAddress", Value: c.SendAddress},
		{Path: "SendInterv", Value: c.SendInterv},
		{Path: "SslEnable", Value: c.SslEnable},
		{Path: "Title", Value: c.Title},
		{Path: "TlsEnable", Value: c.TLSEnable},
		{Path: "UserName", Value: c.UserName},
	})
}

func (g Email) Validate() error {
	return nil
}
