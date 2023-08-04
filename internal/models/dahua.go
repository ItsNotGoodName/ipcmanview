package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type DahuaCamera struct {
	ID        int64
	Address   string `validate:"address"`
	Username  string
	Password  string
	CreatedAt time.Time
}

type DahuaCameraDetail struct {
	ID              int64
	SN              string
	DeviceClass     string
	DeviceType      string
	HardwareVersion string
	MarketArea      string
	ProcessInfo     string
	Vendor          string
}

type DahuaSoftwareVersion struct {
	ID                      int64
	Build                   string
	BuildDate               string
	SecurityBaseLineVersion string
	Version                 string
	WebVersion              string
}

type DahuaCameraLicense struct {
	ID            int64
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

type DahuaScanCamera struct {
	ID           int64
	Seed         int
	Location     *time.Location
	FullComplete bool
	FullCursor   time.Time
	FullEpoch    time.Time
	QuickCursor  time.Time
}

type DahuaScanType string

var (
	DahuaScanTypeFull   = DahuaScanType("full")
	DahuaScanTypeQuick  = DahuaScanType("quick")
	DahuaScanTypeManual = DahuaScanType("manual")
)

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

func (d DahuaScanRange) ScanNull() error {
	return nil
}

func (d DahuaScanRange) ScanBounds() (any, any) {
	return &d.Start, &d.End
}

func (d DahuaScanRange) SetBoundTypes(lower, upper pgtype.BoundType) error {
	return nil
}

type DahuaScanTask struct {
	ID        int64
	CameraID  int64
	ScanRange DahuaScanRange
	Type      DahuaScanType
}

type DahuaScanActiveTask struct {
	ID        int64
	CameraID  int64
	ScanRange DahuaScanRange
	Type      DahuaScanType
	StartedAt time.Time
	Duration  int
	Upserted  int
	Deleted   int
	Percent   float64
}

type DahuaScanCompleteTask struct {
	ID        int64
	CameraID  int64
	ScanRange DahuaScanRange
	Type      DahuaScanType
	StartedAt time.Time
	Duration  int
	Upserted  int
	Deleted   int
	Percent   float64
	Success   bool
	Error     string
}
