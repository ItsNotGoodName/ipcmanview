package main

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

type Shared struct {
	DBPath string `default:"sqlite.db" env:"DB_PATH" help:"Path to SQLite database."`
}

func useDB(path string) (sqlc.DB, error) {
	sqlDB, err := sqlite.New(path)
	if err != nil {
		return sqlc.DB{}, err
	}
	if err := migrations.Migrate(sqlDB); err != nil {
		return sqlc.DB{}, err
	}
	return sqlc.NewDB(sqlite.NewDebugDB(sqlDB)), nil
}
