package main

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

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
