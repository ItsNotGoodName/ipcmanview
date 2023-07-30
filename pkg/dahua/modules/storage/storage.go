package storage

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
)

func GetDeviceAllInfo(ctx context.Context, gen dahua.GenRPC) ([]Storage, error) {
	var object int64
	{
		rpc, err := gen.RPC(ctx)
		if err != nil {
			return []Storage{}, err
		}

		res, err := dahua.Send[any](ctx, rpc.Method("storage.factory.instance"))
		if err != nil {
			return []Storage{}, err
		}

		object = res.Result.Number
	}

	rpc, err := gen.RPC(ctx)
	if err != nil {
		return []Storage{}, err
	}

	res, err := dahua.Send[GetDeviceAllInfoResult](ctx, rpc.Method("storage.getDeviceAllInfo").Object(object))
	if err != nil {
		return []Storage{}, err
	}

	return res.Params.Info, err
}

type Storage struct {
	Name   string          `json:"Name"`
	State  string          `json:"State"`
	Detail []StorageDetail `json:"Detail"`
}

type StorageDetail struct {
	Path       string  `json:"Path"`
	Type       string  `json:"Type"`
	TotalBytes float64 `json:"TotalBytes"`
	UsedBytes  float64 `json:"UsedBytes"`
	IsError    bool    `json:"IsError"`
}

type GetDeviceAllInfoResult struct {
	Info []Storage
}

func (g *GetDeviceAllInfoResult) UnmarshalJSON(data []byte) error {
	{
		res := struct {
			Info []Storage `json:"info"`
		}{}

		if err := json.Unmarshal(data, &res); err == nil {
			g.Info = res.Info
			return nil
		}
	}

	var storages []Storage
	if err := json.Unmarshal(data, &storages); err != nil {
		return err
	}

	g.Info = storages

	return nil
}
