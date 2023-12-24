package configmanager

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Config[T any] struct {
	Table    T
	Original json.RawMessage
}

// func (c Config[T]) Merge() (json.RawMessage, error) {
// 	b, err := json.Marshal(c.Table)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	left := make(map[string]any)
// 	if err := json.Unmarshal(b, &left); err != nil {
// 		return nil, err
// 	}
//
// 	right := make(map[string]any)
// 	if err := json.Unmarshal(c.Original, &right); err != nil {
// 		return nil, err
// 	}
//
// 	for key, value := range left {
// 		right[key] = value
// 	}
//
// 	b, err = json.Marshal(right)
// 	return json.RawMessage(b), err
// }

func GetConfig[T any](ctx context.Context, c dahuarpc.Conn, name string) (Config[T], error) {
	res, err := dahuarpc.Send[struct {
		Table json.RawMessage `json:"table"`
	}](ctx, c, dahuarpc.
		New("configManager.getConfig").
		Params(struct {
			Name string `json:"name"`
		}{
			Name: name,
		}))
	if err != nil {
		return Config[T]{}, err
	}

	var table T
	if err := json.Unmarshal(res.Params.Table, &table); err != nil {
		return Config[T]{}, err
	}

	return Config[T]{
		Original: res.Params.Table,
		Table:    table,
	}, nil
}

func SetConfig[T any](ctx context.Context, c dahuarpc.Conn, name string, config Config[T]) error {
	_, err := dahuarpc.Send[any](ctx, c, dahuarpc.
		New("configManager.setConfig").
		Params(struct {
			Name  string          `json:"name"`
			Table json.RawMessage `json:"table"`
		}{
			Name:  name,
			Table: config.Original,
		}))
	return err
}
