package models

import (
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

type Error struct {
	Error string `json:"error"`
}

type StreamPayload struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	OK      bool    `json:"ok"`
}
