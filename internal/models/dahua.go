package models

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dbgen/postgres/dahua/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type DahuaCamera struct {
	ID        int64
	Address   string `validate:"address"`
	Username  string
	Password  string
	Location  Location
	CreatedAt time.Time
}

type DahuaCameraDetail struct {
	SN              string
	DeviceClass     string
	DeviceType      string
	HardwareVersion string
	MarketArea      string
	ProcessInfo     string
	Vendor          string
}

type DahuaSoftwareVersion struct {
	Build                   string
	BuildDate               string
	SecurityBaseLineVersion string
	Version                 string
	WebVersion              string
}

type DahuaCameraLicense struct {
	AbroadInfo    string
	AllType       bool
	DigitChannel  int
	EffectiveDays int
	EffectiveTime int
	LicenseID     int
	ProductType   string
	Status        int
	Username      string
}

type DahuaScanCursor struct {
	CameraID     int64
	Seed         int
	Location     Location
	FullComplete bool
	FullCursor   time.Time
	FullEpoch    time.Time
	FullEpochEnd time.Time
	QuickCursor  time.Time
}

// DahuaScanRange is INCLUSIVE Start and EXCLUSIVE End.
type DahuaScanRange struct {
	Start time.Time
	End   time.Time
}

func (d DahuaScanRange) IsNull() bool {
	return false
}

func (d DahuaScanRange) BoundTypes() (pgtype.BoundType, pgtype.BoundType) {
	return pgtype.Inclusive, pgtype.Exclusive
}

func (d DahuaScanRange) Bounds() (any, any) {
	return d.Start, d.End
}

func (d *DahuaScanRange) ScanNull() error {
	return nil
}

func (d *DahuaScanRange) ScanBounds() (any, any) {
	return &d.Start, &d.End
}

func (d *DahuaScanRange) SetBoundTypes(lower, upper pgtype.BoundType) error {
	return nil
}

type DahuaScanKind = model.ScanKind

var (
	DahuaScanKindFull   = model.ScanKind_Full
	DahuaScanKindQuick  = model.ScanKind_Quick
	DahuaScanKindManual = model.ScanKind_Manual
)

type DahuaScanQueueTask struct {
	ID       int64
	CameraID int64
	Kind     DahuaScanKind
	Range    DahuaScanRange
}

type DahuaScanActiveTask struct {
	CameraID  int64
	Kind      DahuaScanKind
	Range     DahuaScanRange
	Cursor    time.Time
	StartedAt time.Time
	Deleted   int
	Upserted  int
	Percent   float64
}

func (q DahuaScanActiveTask) NewProgress() DahuaScanActiveProgress {
	return DahuaScanActiveProgress{
		CameraID: q.CameraID,
	}
}

type DahuaScanActiveProgress struct {
	CameraID int64
	Upserted int
	Deleted  int
	Percent  float64
	Cursor   time.Time
}

type DahuaScanCompleteTask struct {
	ID        int64
	CameraID  int64
	Kind      DahuaScanKind
	Range     DahuaScanRange
	Cursor    time.Time
	StartedAt time.Time
	Duration  int
	Upserted  int
	Deleted   int
	Percent   float64
	Error     string
}

type DahuaCameraEvent struct {
	ID            int64
	CameraID      int64
	ContentType   string
	ContentLength int
	Code          string
	Action        string
	Index         int
	Data          []byte
	CreatedAt     time.Time
}
