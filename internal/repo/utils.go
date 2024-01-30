package repo

import (
	"database/sql"
	"errors"
)

var ErrNotFound = sql.ErrNoRows

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func Coalasce(s ...*string) string {
	for _, v := range s {
		if v != nil {
			return *v
		}
	}
	return ""
}
