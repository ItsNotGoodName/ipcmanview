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

type User struct {
	ID       int64
	Email    string `validate:"required,lte=128,email,excludes= "`
	Username string `validate:"gte=3,lte=64,excludes=@,excludes= "`
	Password string `validate:"gte=8"`
}

type AuthSession struct {
	SessionID int64
	UserID    int64
	Username  string
	Session   string
	Admin     bool
}
