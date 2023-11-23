package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Location struct {
	*time.Location
}

func (l *Location) MarshalJSON() ([]byte, error) {
	return []byte(l.Location.String()), nil
}

func (l *Location) UnmarshalJSON(data []byte) error {
	loc, err := time.LoadLocation(string(data))
	if err != nil {
		return err
	}
	*l = Location{loc}
	return nil
}

func (dst *Location) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case string:
		loc, err := time.LoadLocation(string(src))
		if err != nil {
			return err
		}
		*dst = Location{loc}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Location) Value() (driver.Value, error) {
	return src.Location.String(), nil
}

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

type Error struct {
	Error string `json:"error"`
}

type StreamPayload struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	OK      bool    `json:"ok"`
}
