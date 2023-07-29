package magicbox

import (
	"encoding/json"

	"github.com/ItsNotGoodName/pkg/dahua"
)

func Reboot(gen dahua.Generator) (bool, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return false, err
	}

	res, err := dahua.Send[any](rpc.Method("magicBox.reboot"))

	return res.Result.Bool(), err
}

func NeedReboot(gen dahua.Generator) (int, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return 0, err
	}

	res, err := dahua.Send[struct {
		NeedReboot int `json:"needReboot"`
	}](rpc.Method("magicBox.needReboot"))

	return res.Params.NeedReboot, err
}

func GetSerialNo(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		SN string `json:"sn"`
	}](rpc.Method("magicBox.getSerialNo"))

	return res.Params.SN, err
}

func GetDeviceType(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Type string `json:"type"`
	}](rpc.Method("magicBox.getDeviceType"))

	return res.Params.Type, err
}

type GetMemoryInfoResult struct {
	Free  int64 `json:"free"`
	Total int64 `json:"total"`
}

func (g *GetMemoryInfoResult) UnmarshalJSON(data []byte) error {
	var res struct {
		Free  float64 `json:"free"`
		Total float64 `json:"total"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	g.Free = int64(res.Free)
	g.Total = int64(res.Total)

	return nil
}

func GetMemoryInfo(gen dahua.Generator) (GetMemoryInfoResult, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return GetMemoryInfoResult{}, err
	}

	res, err := dahua.Send[GetMemoryInfoResult](rpc.Method("magicBox.getMemoryInfo"))

	return res.Params, err
}

func GetCPUUsage(gen dahua.Generator) (int, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return 0, err
	}

	res, err := dahua.Send[struct {
		Usage int `json:"usage"`
	}](rpc.Method("magicBox.getCPUUsage").Params(dahua.JSON{"index": 0}))

	return res.Params.Usage, err
}

func GetDeviceClass(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Type string `json:"type"`
	}](rpc.Method("magicBox.getDeviceClass"))

	return res.Params.Type, err
}

func GetProcessInfo(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Info string `json:"info"`
	}](rpc.Method("magicBox.getProcessInfo"))

	return res.Params.Info, err
}

func GetHardwareVersion(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Version string `json:"version"`
	}](rpc.Method("magicBox.getHardwareVersion"))

	return res.Params.Version, err
}

func GetVendor(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Vendor string `json:"Vendor"`
	}](rpc.Method("magicBox.getVendor"))

	return res.Params.Vendor, err
}

type GetSoftwareVersionResult struct {
	Build                   string `json:"Build"`
	BuildDate               string `json:"BuildDate"`
	SecurityBaseLineVersion string `json:"SecurityBaseLineVersion"`
	Version                 string `json:"Version"`
	WebVersion              string `json:"WebVersion"`
}

func GetSoftwareVersion(gen dahua.Generator) (GetSoftwareVersionResult, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return GetSoftwareVersionResult{}, err
	}

	res, err := dahua.Send[struct {
		Version GetSoftwareVersionResult `json:"version"`
	}](rpc.Method("magicBox.getSoftwareVersion"))
	return res.Params.Version, err
}

func GetMarketArea(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		AbroadInfo string `json:"AbroadInfo"`
	}](rpc.Method("magicBox.getMarketArea"))

	return res.Params.AbroadInfo, err
}

type GetUptimeResult struct {
	Last  int64 `json:"last"`
	Total int64 `json:"total"`
}

func GetUpTime(gen dahua.Generator) (GetMemoryInfoResult, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return GetMemoryInfoResult{}, err
	}

	res, err := dahua.Send[struct {
		Info GetMemoryInfoResult `json:"info"`
	}](rpc.Method("magicBox.getUpTime"))

	return res.Params.Info, err
}
