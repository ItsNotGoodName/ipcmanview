package models

import "time"

// TimeRange is INCLUSIVE Start and EXCLUSIVE End.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

type Error struct {
	Error string `json:"error"`
}

type StreamPayload struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	OK      bool    `json:"ok"`
}
