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
	value    DahuaCamera
	Address  bool
	Username bool
	Password bool
}

func NewDahuaCameraUpdate(id int64) *DahuaCameraUpdate {
	return &DahuaCameraUpdate{value: DahuaCamera{ID: id}}
}

func (d *DahuaCameraUpdate) UpdateAddress(addresss string) *DahuaCameraUpdate {
	d.value.Address = addresss
	d.Address = true
	return d
}

func (d *DahuaCameraUpdate) UpdateUsername(username string) *DahuaCameraUpdate {
	d.value.Username = username
	d.Username = true
	return d
}

func (d *DahuaCameraUpdate) UpdatePassword(password string) *DahuaCameraUpdate {
	d.value.Password = password
	d.Password = true
	return d
}

func (d *DahuaCameraUpdate) Value() (DahuaCamera, error) {
	if err := validate.Struct(&d.value); err != nil {
		return DahuaCamera{}, err
	}

	return d.value, nil
}
