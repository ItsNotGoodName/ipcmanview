package repo

import (
	"database/sql"
	"errors"
)

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func ErrorToNullString(err error) sql.NullString {
	if err == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: err.Error(),
		Valid:  true,
	}
}
