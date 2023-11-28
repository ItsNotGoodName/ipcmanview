package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice struct {
	Slice []string
}

func (dst *StringSlice) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case string:
		return json.Unmarshal([]byte(src), &dst.Slice)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src StringSlice) Value() (driver.Value, error) {
	b, err := json.Marshal(src.Slice)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
