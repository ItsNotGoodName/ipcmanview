//go:build !cgo

package sqlite

import (
	"database/sql"
	"errors"

	"modernc.org/sqlite"
)

func connect(dbPath string) (*sql.DB, error) {
	pragmas := "?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)"
	db, err := sql.Open("sqlite", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AsError(err error) (Error, bool) {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return Error{
			Code: sqliteErr.Code(),
			Msg:  sqliteErr.Error(),
		}, true
	}
	return Error{}, false
}
