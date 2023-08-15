package magicbox

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func Reboot(ctx context.Context, c dahuarpc.Client) (bool, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahuarpc.Send[any](ctx, rpc.Method("magicBox.reboot"))

	return res.Result.Bool(), err
}

func NeedReboot(ctx context.Context, c dahuarpc.Client) (int, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahuarpc.Send[struct {
		NeedReboot int `json:"needReboot"`
	}](ctx, rpc.Method("magicBox.needReboot"))

	return res.Params.NeedReboot, err
}

func GetSerialNo(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		SN string `json:"sn"`
	}](ctx, rpc.Method("magicBox.getSerialNo"))

	return res.Params.SN, err
}

func GetDeviceType(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Type string `json:"type"`
	}](ctx, rpc.Method("magicBox.getDeviceType"))

	return res.Params.Type, err
}

func GetMemoryInfo(ctx context.Context, c dahuarpc.Client) (MemoryInfo, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return MemoryInfo{}, err
	}

	res, err := dahuarpc.Send[MemoryInfo](ctx, rpc.Method("magicBox.getMemoryInfo"))

	return res.Params, err
}

type MemoryInfo struct {
	Free  dahuarpc.Integer `json:"free"`
	Total dahuarpc.Integer `json:"total"`
}

func GetCPUUsage(ctx context.Context, c dahuarpc.Client) (int, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahuarpc.Send[struct {
		Usage int `json:"usage"`
	}](ctx, rpc.Method("magicBox.getCPUUsage").Params(dahuarpc.JSON{"index": 0}))

	return res.Params.Usage, err
}

func GetDeviceClass(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Type string `json:"type"`
	}](ctx, rpc.Method("magicBox.getDeviceClass"))

	return res.Params.Type, err
}

func GetProcessInfo(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Info string `json:"info"`
	}](ctx, rpc.Method("magicBox.getProcessInfo"))

	return res.Params.Info, err
}

func GetHardwareVersion(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Version string `json:"version"`
	}](ctx, rpc.Method("magicBox.getHardwareVersion"))

	return res.Params.Version, err
}

func GetVendor(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Vendor string `json:"Vendor"`
	}](ctx, rpc.Method("magicBox.getVendor"))

	return res.Params.Vendor, err
}

func GetSoftwareVersion(ctx context.Context, c dahuarpc.Client) (SoftwareVersion, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return SoftwareVersion{}, err
	}

	res, err := dahuarpc.Send[struct {
		Version SoftwareVersion `json:"version"`
	}](ctx, rpc.Method("magicBox.getSoftwareVersion"))
	return res.Params.Version, err
}

type SoftwareVersion struct {
	Build                   string `json:"Build"`
	BuildDate               string `json:"BuildDate"`
	SecurityBaseLineVersion string `json:"SecurityBaseLineVersion"`
	Version                 string `json:"Version"`
	WebVersion              string `json:"WebVersion"`
}

func GetMarketArea(ctx context.Context, c dahuarpc.Client) (string, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		AbroadInfo string `json:"AbroadInfo"`
	}](ctx, rpc.Method("magicBox.getMarketArea"))

	return res.Params.AbroadInfo, err
}

func GetUpTime(ctx context.Context, c dahuarpc.Client) (UpTime, error) {
	rpc, err := c.RPC(ctx)
	if err != nil {
		return UpTime{}, err
	}

	res, err := dahuarpc.Send[struct {
		Info UpTime `json:"info"`
	}](ctx, rpc.Method("magicBox.getUpTime"))

	return res.Params.Info, err
}

type UpTime struct {
	Last  int64 `json:"last"`
	Total int64 `json:"total"`
}
