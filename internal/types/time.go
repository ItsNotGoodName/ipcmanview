package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func NewTime(t time.Time) Time {
	return Time{Time: t}
}

// Time will always UTC.
type Time struct {
	time.Time
}

func (src Time) Value() (driver.Value, error) {
	return src.Time.UTC(), nil
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

func NewNullTime(t time.Time) NullTime {
	return NullTime{
		Time:  t,
		Valid: true,
	}
}

// NullTime will always UTC.
// TODO: replace this with sql.Null[T] type when Go 1.22 comes out
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (dst *NullTime) Scan(src any) error {
	if src == nil {
		dst.Time, dst.Valid = time.Time{}, false
		return nil
	}
	dst.Valid = true

	switch src := src.(type) {
	case time.Time:
		dst.Time = src.UTC()
		return nil
	}

	return fmt.Errorf("cannot scan %T", dst)
}

func (src NullTime) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Time.UTC(), nil
}
