//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type CameraFiles struct {
	ID        int32 `sql:"primary_key"`
	CameraID  int32
	FilePath  string
	Kind      string
	Size      int32
	StartTime time.Time
	EndTime   time.Time
	Duration  int32
	ScannedAt time.Time
	Events    string
}
