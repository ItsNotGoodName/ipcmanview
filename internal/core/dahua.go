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

type DahuaCameraCreate struct {
	Address  string
	Username string
	Password string
}

func NewDahuaCamera(r DahuaCameraCreate) (DahuaCamera, error) {
	res := DahuaCamera{
		Address:  r.Address,
		Username: r.Username,
		Password: r.Password,
	}

	err := validate.Struct(&res)

	return res, err
}

type DahuaCameraUpdate struct {
	DahuaCamera
	Address  bool
	Username bool
	Password bool
}

func NewDahuaCameraUpdate(id int64) DahuaCameraUpdate {
	return DahuaCameraUpdate{DahuaCamera: DahuaCamera{ID: id}}
}

func (d *DahuaCameraUpdate) UpdateAddress(addresss string) error {
	d.DahuaCamera.Address = addresss
	d.Address = true
	return validate.StructPartial(d.DahuaCamera, "Address")
}

func (d *DahuaCameraUpdate) UpdateUsername(username string) error {
	d.DahuaCamera.Username = username
	d.Username = true
	return validate.StructPartial(d.DahuaCamera, "Username")
}

func (d *DahuaCameraUpdate) UpdatePassword(password string) error {
	d.DahuaCamera.Password = password
	d.Password = true
	return validate.StructPartial(d.DahuaCamera, "Password")
}
