package models

import "time"

// TimeRange is INCLUSIVE Start and EXCLUSIVE End.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (t TimeRange) Null() bool {
	return t.Start.IsZero() && t.End.IsZero()
}
