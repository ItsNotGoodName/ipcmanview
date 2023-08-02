package core

import (
	"time"
)

type DahuaCamera struct {
	ID        int64
	Address   string `validate:"address"`
	Username  string
	Password  string
	CreatedAt time.Time
}

func (dc DahuaCamera) Equal(cam DahuaCamera) bool {
	return dc.Address == cam.Address &&
		dc.Username == cam.Username &&
		dc.Password == cam.Username
}

func (dc DahuaCamera) Validate() (DahuaCamera, error) {
	err := validate.Struct(&dc)
	if err != nil {
		return DahuaCamera{}, err
	}

	return dc, nil
}

type DahuaCameraCreate struct {
	Address  string
	Username string
	Password string
}

func NewDahuaCamera(r DahuaCameraCreate) (DahuaCamera, error) {
	dc := DahuaCamera{
		Address:  r.Address,
		Username: r.Username,
		Password: r.Password,
	}
	return dc, validate.Struct(&dc)
}

type DahuaCameraUpdate struct {
	value    DahuaCamera
	Address  bool
	Username bool
	Password bool
}

func NewDahuaCameraUpdate(id int64) *DahuaCameraUpdate {
	return &DahuaCameraUpdate{value: DahuaCamera{ID: id}}
}

func (d *DahuaCameraUpdate) AddressUpdate(addresss string) *DahuaCameraUpdate {
	d.value.Address = addresss
	d.Address = true
	return d
}

func (d *DahuaCameraUpdate) UsernameUpdate(username string) *DahuaCameraUpdate {
	d.value.Username = username
	d.Username = true
	return d
}

func (d *DahuaCameraUpdate) PasswordUpdate(password string) *DahuaCameraUpdate {
	d.value.Password = password
	d.Password = true
	return d
}

func (d *DahuaCameraUpdate) Value() (DahuaCamera, error) {
	return d.value.Validate()
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
