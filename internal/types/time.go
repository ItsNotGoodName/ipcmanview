package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func NewTime(t time.Time) Time {
	return Time{Time: t}
}

type Time struct {
	time.Time
}

func (d Time) Value() (driver.Value, error) {
	return d.Time.UTC(), nil
}

func (dst *Time) Scan(src any) error {
	if dst == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case time.Time:
		dst.Time = src.UTC()
		return nil
	}

	return fmt.Errorf("cannot scan %T", dst)
}
