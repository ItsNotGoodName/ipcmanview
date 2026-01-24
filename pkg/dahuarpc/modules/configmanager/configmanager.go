package configmanager

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type ConfigData interface {
	// Merge should merge its own data with the input js.
	// The return should be the merged data.
	Merge(js string) (string, error)
	// Validate its own data.
	Validate() error
}

type ConfigTable[T ConfigData] struct {
	Data T
	raw  string
}

func (ct *ConfigTable[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &ct.Data); err != nil {
		return err
	}

	err := ct.Data.Validate()
	if err != nil {
		return err
	}

	ct.raw = string(data)

	return nil
}

func (ct *ConfigTable[T]) MarshalJSON() ([]byte, error) {
	b, err := ct.Data.Merge(string(ct.raw))
	if err != nil {
		return nil, err
	}
	return []byte(b), err
}

type Config[T ConfigData] struct {
	name   string
	Tables ConfigTable[T]
}

func GetConfig[T ConfigData](ctx context.Context, c dahuarpc.Conn, name string) (Config[T], error) {
	rb := dahuarpc.
		New("configManager.getConfig").
		Params(struct {
			Name string `json:"name"`
		}{
			Name: name,
		})

	res, err := dahuarpc.Send[struct {
		Table ConfigTable[T] `json:"table"`
	}](ctx, c, rb)
	if err != nil {
		return Config[T]{}, err
	}

	return Config[T]{
		name:   name,
		Tables: res.Params.Table,
	}, nil
}

func SetConfig[T ConfigData](ctx context.Context, c dahuarpc.Conn, config Config[T]) error {
	_, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("configManager.setConfig").
		Params(struct {
			Name  string         `json:"name"`
			Table ConfigTable[T] `json:"table"`
		}{
			Name:  config.name,
			Table: config.Tables,
		}))
	return err
}

type ConfigArray[T ConfigData] struct {
	name   string
	Tables []ConfigTable[T]
}

func GetConfigArray[T ConfigData](ctx context.Context, c dahuarpc.Conn, name string) (ConfigArray[T], error) {
	rb := dahuarpc.
		New("configManager.getConfig").
		Params(struct {
			Name string `json:"name"`
		}{
			Name: name,
		})

	res, err := dahuarpc.Send[struct {
		Table []ConfigTable[T] `json:"table"`
	}](ctx, c, rb)
	if err != nil {
		return ConfigArray[T]{}, err
	}

	return ConfigArray[T]{
		name:   name,
		Tables: res.Params.Table,
	}, nil
}

func SetConfigArray[T ConfigData](ctx context.Context, c dahuarpc.Conn, config ConfigArray[T]) error {
	_, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("configManager.setConfig").
		Params(struct {
			Name  string           `json:"name"`
			Table []ConfigTable[T] `json:"table"`
		}{
			Name:  config.name,
			Table: config.Tables,
		}))
	return err
}
