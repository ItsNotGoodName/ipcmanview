package dahua

import "database/sql"

func errorToNullString(err error) sql.NullString {
	if err == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: err.Error(),
		Valid:  true,
	}
}
