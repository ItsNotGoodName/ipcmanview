//go:build cgo

package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		PRAGMA busy_timeout       = 5000;
		PRAGMA journal_mode       = WAL;
		PRAGMA foreign_keys       = ON;
	`, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}
