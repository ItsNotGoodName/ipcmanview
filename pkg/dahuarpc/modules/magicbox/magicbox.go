package magicbox

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

func Reboot(ctx context.Context, c dahuarpc.Conn) (bool, error) {
	res, err := dahuarpc.Send[any](ctx, c, dahuarpc.New("magicBox.reboot"))

	return res.Result.Bool(), err
}

func NeedReboot(ctx context.Context, c dahuarpc.Conn) (int, error) {
	res, err := dahuarpc.Send[struct {
		NeedReboot int `json:"needReboot"`
	}](ctx, c, dahuarpc.New("magicBox.needReboot"))

	return res.Params.NeedReboot, err
}

func GetSerialNo(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		SN string `json:"sn"`
	}](ctx, c, dahuarpc.New("magicBox.getSerialNo"))

	return res.Params.SN, err
}

func GetDeviceType(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Type string `json:"type"`
	}](ctx, c, dahuarpc.New("magicBox.getDeviceType"))

	return res.Params.Type, err
}

func GetMemoryInfo(ctx context.Context, c dahuarpc.Conn) (MemoryInfo, error) {
	res, err := dahuarpc.Send[MemoryInfo](ctx, c, dahuarpc.New("magicBox.getMemoryInfo"))

	return res.Params, err
}

type MemoryInfo struct {
	Free  dahuarpc.Integer `json:"free"`
	Total dahuarpc.Integer `json:"total"`
}

func GetCPUUsage(ctx context.Context, c dahuarpc.Conn) (int, error) {
	res, err := dahuarpc.Send[struct {
		Usage int `json:"usage"`
	}](ctx, c, dahuarpc.
		New("magicBox.getCPUUsage").
		Params(struct {
			Index int `json:"index"`
		}{
			Index: 0,
		}))

	return res.Params.Usage, err
}

func GetDeviceClass(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Type string `json:"type"`
	}](ctx, c, dahuarpc.New("magicBox.getDeviceClass"))

	return res.Params.Type, err
}

func GetProcessInfo(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Info string `json:"info"`
	}](ctx, c, dahuarpc.New("magicBox.getProcessInfo"))

	return res.Params.Info, err
}

func GetHardwareVersion(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Version string `json:"version"`
	}](ctx, c, dahuarpc.New("magicBox.getHardwareVersion"))

	return res.Params.Version, err
}

func GetVendor(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Vendor string `json:"Vendor"`
	}](ctx, c, dahuarpc.New("magicBox.getVendor"))

	return res.Params.Vendor, err
}

func GetSoftwareVersion(ctx context.Context, c dahuarpc.Conn) (SoftwareVersion, error) {
	res, err := dahuarpc.Send[struct {
		Version SoftwareVersion `json:"version"`
	}](ctx, c, dahuarpc.New("magicBox.getSoftwareVersion"))
	return res.Params.Version, err
}

type SoftwareVersion struct {
	Build                   string `json:"Build"`
	BuildDate               string `json:"BuildDate"`
	SecurityBaseLineVersion string `json:"SecurityBaseLineVersion"`
	Version                 string `json:"Version"`
	WebVersion              string `json:"WebVersion"`
}

func GetMarketArea(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		AbroadInfo string `json:"AbroadInfo"`
	}](ctx, c, dahuarpc.New("magicBox.getMarketArea"))

	return res.Params.AbroadInfo, err
}

func GetUpTime(ctx context.Context, c dahuarpc.Conn) (UpTime, error) {
	res, err := dahuarpc.Send[struct {
		Info UpTime `json:"info"`
	}](ctx, c, dahuarpc.New("magicBox.getUpTime"))

	return res.Params.Info, err
}

type UpTime struct {
	Last  int64 `json:"last"`
	Total int64 `json:"total"`
}

func GetMachineName(ctx context.Context, c dahuarpc.Conn) (string, error) {
	res, err := dahuarpc.Send[struct {
		Name string `json:"name"`
	}](ctx, c, dahuarpc.New("magicBox.getMachineName"))
	return res.Params.Name, err
}

func ListMethod(ctx context.Context, c dahuarpc.Conn) ([]string, error) {
	res, err := dahuarpc.Send[struct {
		Method []string `json:"method"`
	}](ctx, c, dahuarpc.New("magicBox.listMethod"))
	return res.Params.Method, err
}
