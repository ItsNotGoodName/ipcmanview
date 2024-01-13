// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package repo

import (
	"database/sql"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

type DahuaAferoFile struct {
	ID                int64
	FileID            sql.NullInt64
	FileThumbnailID   sql.NullInt64
	EmailAttachmentID sql.NullInt64
	Name              string
	CreatedAt         types.Time
	Ready             bool
}

type DahuaDevice struct {
	ID        int64
	Name      string
	Ip        string
	Url       types.URL
	Username  string
	Password  string
	Location  types.Location
	Feature   models.DahuaFeature
	CreatedAt types.Time
	UpdatedAt types.Time
}

type DahuaEmailAttachment struct {
	ID        int64
	MessageID int64
	FileName  string
}

type DahuaEmailMessage struct {
	ID                int64
	DeviceID          int64
	Date              types.Time
	From              string
	To                types.StringSlice
	Subject           string
	Text              string
	AlarmEvent        string
	AlarmInputChannel int64
	AlarmName         string
	CreatedAt         types.Time
}

type DahuaEvent struct {
	ID        int64
	DeviceID  int64
	Code      string
	Action    string
	Index     int64
	Data      json.RawMessage
	CreatedAt types.Time
}

type DahuaEventDeviceRule struct {
	DeviceID   int64
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
}

type DahuaEventRule struct {
	ID         int64
	Code       string
	IgnoreDb   bool
	IgnoreLive bool
	IgnoreMqtt bool
}

type DahuaEventWorkerState struct {
	ID        int64
	DeviceID  int64
	State     models.DahuaEventWorkerState
	Error     sql.NullString
	CreatedAt types.Time
}

type DahuaFile struct {
	ID          int64
	DeviceID    int64
	Channel     int64
	StartTime   types.Time
	EndTime     types.Time
	Length      int64
	Type        string
	FilePath    string
	Duration    int64
	Disk        int64
	VideoStream string
	Flags       types.StringSlice
	Events      types.StringSlice
	Cluster     int64
	Partition   int64
	PicIndex    int64
	Repeat      int64
	WorkDir     string
	WorkDirSn   bool
	UpdatedAt   types.Time
	Storage     models.Storage
}

type DahuaFileCursor struct {
	DeviceID     int64
	QuickCursor  types.Time
	FullCursor   types.Time
	FullEpoch    types.Time
	FullComplete bool
	Scan         bool
	ScanPercent  float64
	ScanType     models.DahuaScanType
}

type DahuaFileThumbnail struct {
	ID     int64
	FileID int64
	Width  int64
	Height int64
}

type DahuaSeed struct {
	Seed     int64
	DeviceID sql.NullInt64
}

type DahuaStorageDestination struct {
	ID              int64
	Name            string
	Storage         models.Storage
	ServerAddress   string
	Port            int64
	Username        string
	Password        string
	RemoteDirectory string
}

type DahuaStream struct {
	ID           int64
	Internal     bool
	DeviceID     int64
	Channel      int64
	Subtype      int64
	Name         string
	MediamtxPath string
}

type Setting struct {
	SiteName        string
	DefaultLocation types.Location
}
