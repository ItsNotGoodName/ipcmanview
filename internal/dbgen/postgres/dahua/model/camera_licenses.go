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

type CameraLicenses struct {
	CameraID      int32
	AbroadInfo    string
	AllType       bool
	DigitChannel  int32
	EffectiveDays int32
	EffectiveTime int32
	EffectiveAt   time.Time
	LicenseID     int32
	ProductType   string
	Status        int32
	Username      string
}