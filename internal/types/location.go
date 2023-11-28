package types

import (
	"database/sql/driver"
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
