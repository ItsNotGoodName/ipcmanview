package core

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound  = fmt.Errorf("not found")
	ErrForbidden = fmt.Errorf("forbidden")
)

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, sql.ErrNoRows)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}
