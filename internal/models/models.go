package models

import "time"

type EventAction string

// TimeRange is INCLUSIVE Start and EXCLUSIVE End.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (t TimeRange) Null() bool {
	return t.Start.IsZero() && t.End.IsZero()
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Storage string

const (
	StorageLocal Storage = "local"
	StorageSFTP  Storage = "sftp"
	StorageFTP   Storage = "ftp"
	// StorageNFS   Storage = "nfs"
	// StorageSMB   Storage = "smb"
)
