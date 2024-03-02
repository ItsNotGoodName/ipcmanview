package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var timeFormats = []string{
	"2006-01-02 15:04:05.000000",
	"2006-01-02 15:04:05.000000 +0000 UTC",
	"2006-01-02 15:04:05",
}

func NewTime(t time.Time) Time {
	return Time{Time: t}
}

// Time will always UTC.
type Time struct {
	time.Time
}

func (dst *Time) Scan(src any) error {
	switch src := src.(type) {
	case time.Time:
		dst.Time = src.UTC()
		return nil
	case string:
		for _, f := range timeFormats {
			t, err := time.ParseInLocation(f, src, time.UTC)
			if err != nil {
				continue
			}
			dst.Time = t
			return nil
		}

		return fmt.Errorf("parsing time %s", src)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Time) Value() (driver.Value, error) {
	return src.Time.UTC().Format(timeFormats[0]), nil
}

// NullTime will always UTC.
type NullTime struct {
	Time
	Valid bool // Valid is true if Time is not NULL
}

func (dst *NullTime) Scan(src any) error {
	if src == nil {
		dst.Time, dst.Valid = Time{}, false
		return nil
	}

	t := &Time{}
	err := t.Scan(src)
	if err != nil {
		return err
	}
	dst.Time = *t
	dst.Valid = true

	return nil
}

func (src NullTime) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return src.Time.Value()
}
