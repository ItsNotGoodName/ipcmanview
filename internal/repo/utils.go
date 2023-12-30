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

func NilStringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: true,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func Coalasce(s ...*string) string {
	for _, v := range s {
		if v != nil {
			return *v
		}
	}
	return ""
}
