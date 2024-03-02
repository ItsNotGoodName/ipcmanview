package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	if src == nil {
		dst.RawMessage = []byte{}
		return nil
	}

	switch src := src.(type) {
	case string:
		dst.RawMessage = []byte(src)
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src JSON) Value() (driver.Value, error) {
	if len(src.RawMessage) == 0 {
		return nil, nil
	}
	return string(src.RawMessage), nil
}
