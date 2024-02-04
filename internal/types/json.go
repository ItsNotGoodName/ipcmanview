package types

import (
	"database/sql/driver"
	"encoding/json"
)

func NewJSON(data json.RawMessage) JSON {
	return JSON{
		RawMessage: data,
	}
}

type JSON struct {
	json.RawMessage
}

func (dst *JSON) Scan(src any) error {
	switch src := src.(type) {
	case string:
		dst.RawMessage = []byte(src)
	default:
		b, err := json.Marshal(src)
		if err != nil {
			return err
		}
		dst.RawMessage = b
	}
	return nil
}

func (src JSON) Value() (driver.Value, error) {
	if len(src.RawMessage) == 0 {
		return nil, nil
	}
	return string(src.RawMessage), nil
}
