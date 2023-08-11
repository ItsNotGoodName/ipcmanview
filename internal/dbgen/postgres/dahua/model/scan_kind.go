//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type ScanKind string

const (
	ScanKind_Full   ScanKind = "full"
	ScanKind_Quick  ScanKind = "quick"
	ScanKind_Manual ScanKind = "manual"
)

func (e *ScanKind) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "full":
		*e = ScanKind_Full
	case "quick":
		*e = ScanKind_Quick
	case "manual":
		*e = ScanKind_Manual
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ScanKind enum")
	}

	return nil
}

func (e ScanKind) String() string {
	return string(e)
}