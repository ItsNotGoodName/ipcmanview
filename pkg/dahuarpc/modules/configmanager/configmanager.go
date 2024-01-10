package configmanager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type ConfigData interface {
	Merge(js string) (string, error)
	Validate() error
}

type ConfigTable[T ConfigData] struct {
	Data T
	JSON json.RawMessage
}

type Config[T ConfigData] struct {
	name   string
	array  bool
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
	for _, table := range tables {
		var data T
		if err := json.Unmarshal(table, &data); err != nil {
			return Config[T]{}, err
		}

		err := data.Validate()
		if err != nil {
			return Config[T]{}, err
		}

		configTables = append(configTables, ConfigTable[T]{
			Data: data,
			JSON: table,
		})
	}

	return Config[T]{
		name:   name,
		array:  array,
		Tables: configTables,
	}, nil
}

func SetConfig[T ConfigData](ctx context.Context, c dahuarpc.Conn, config Config[T]) error {
	table, err := config.table()
	if err != nil {
		return err
	}

	_, err = dahuarpc.Send[any](ctx, c, dahuarpc.
		New("configManager.setConfig").
		Params(struct {
			Name  string          `json:"name"`
			Table json.RawMessage `json:"table"`
		}{
			Name:  config.name,
			Table: table,
		}))
	return err
}

func (c Config[T]) table() (json.RawMessage, error) {
	if len(c.Tables) == 0 {
		return nil, fmt.Errorf("no tables")
	}

	var tables []json.RawMessage
	for _, table := range c.Tables {
		js, err := table.Data.Merge(string(table.JSON))
		if err != nil {
			return nil, err
		}

		tables = append(tables, json.RawMessage(js))
	}

	if !c.array {
		return tables[0], nil
	}

	b, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(b), nil
}
