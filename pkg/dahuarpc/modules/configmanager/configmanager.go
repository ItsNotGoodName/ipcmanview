package configmanager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type ConfigData interface {
	Validate() error
}

type ConfigTable[T ConfigData] struct {
	Data     T
	Original json.RawMessage
}

type Config[T ConfigData] struct {
	Name   string
	Array  bool
	Tables []ConfigTable[T]
}

func GetConfig[T ConfigData](ctx context.Context, c dahuarpc.Conn, name string, array bool) (Config[T], error) {
	rb := dahuarpc.
		New("configManager.getConfig").
		Params(struct {
			Name string `json:"name"`
		}{
			Name: name,
		})

	var tables []json.RawMessage
	if array {
		res, err := dahuarpc.Send[struct {
			Table []json.RawMessage `json:"table"`
		}](ctx, c, rb)
		if err != nil {
			return Config[T]{}, err
		}

		tables = res.Params.Table
	} else {
		res, err := dahuarpc.Send[struct {
			Table json.RawMessage `json:"table"`
		}](ctx, c, rb)
		if err != nil {
			return Config[T]{}, err
		}

		tables = append(tables, res.Params.Table)
	}
	if len(tables) == 0 {
		return Config[T]{}, fmt.Errorf("no tables")
	}

	var configTables []ConfigTable[T]
	for _, t := range tables {
		var data T
		if err := json.Unmarshal(t, &data); err != nil {
			return Config[T]{}, err
		}

		err := data.Validate()
		if err != nil {
			return Config[T]{}, err
		}

		configTables = append(configTables, ConfigTable[T]{
			Data:     data,
			Original: t,
		})
	}

	return Config[T]{
		Name:   name,
		Array:  array,
		Tables: configTables,
	}, nil
}

func SetConfig[T ConfigData](ctx context.Context, c dahuarpc.Conn, config Config[T]) error {
	table, err := config.merge()
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, c, dahuarpc.
		New("configManager.setConfig").
		Params(struct {
			Name  string          `json:"name"`
			Table json.RawMessage `json:"table"`
		}{
			Name:  config.Name,
			Table: table,
		}))
	return err
}

// merge does a shallow merge on all the tables and returns a JSON object or a JSON array.
func (c Config[T]) merge() (json.RawMessage, error) {
	if len(c.Tables) == 0 {
		return nil, fmt.Errorf("no tables")
	}

	var tables []json.RawMessage
	for _, table := range c.Tables {
		b, err := json.Marshal(table.Data)
		if err != nil {
			return nil, err
		}
		left := make(map[string]any)
		if err := json.Unmarshal(b, &left); err != nil {
			return nil, err
		}

		right := make(map[string]any)
		if err := json.Unmarshal(table.Original, &right); err != nil {
			return nil, err
		}

		for key, value := range left {
			if _, ok := right[key]; !ok {
				continue
			}
			right[key] = value
		}

		b, err = json.Marshal(right)
		if err != nil {
			return nil, err
		}

		tables = append(tables, json.RawMessage(b))
	}

	if !c.Array {
		return tables[0], nil
	}

	b, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(b), nil
}
