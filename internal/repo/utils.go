package repo

import (
	"database/sql"
	"errors"
)

var ErrNotFound = sql.ErrNoRows

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
