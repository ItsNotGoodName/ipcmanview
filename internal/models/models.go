package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Location struct {
	*time.Location
}

func (dst *Location) Scan(src any) error {
	if src == nil {
		*dst = Location{time.Local}
		return nil
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
