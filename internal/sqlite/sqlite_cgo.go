//go:build cgo

package sqlite

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func connect(dbPath string) (*sql.DB, error) {
	pragmas := "?_busy_timeout=10000&_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=on"
	db, err := sql.Open("sqlite3", dbPath+pragmas)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AsError(err error) (Error, bool) {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return Error{
			Code: int(sqliteErr.ExtendedCode),
			Msg:  sqliteErr.Error(),
		}, true
	}
	return Error{}, false
}
