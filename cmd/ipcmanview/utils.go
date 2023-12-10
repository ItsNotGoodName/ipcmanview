package main

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"golang.org/x/net/context"
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

type SharedCameras struct {
	ID  []int64 `help:"Run on camera by ID."`
	All bool    `help:"Run on all cameras."`
}

func (c SharedCameras) useCameras(ctx context.Context, db sqlc.DB) ([]models.DahuaCamera, error) {
	var cameras []models.DahuaCamera
	if c.All {
		dbCameras, err := db.ListDahuaCamera(ctx)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range sqlc.ConvertListDahuaCameraRow(dbCameras) {
			cameras = append(cameras, dbCamera)
		}
	} else {
		dbCameras, err := db.ListDahuaCameraByIDs(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range sqlc.ConvertListDahuaCameraByIDsRow(dbCameras) {
			cameras = append(cameras, dbCamera)
		}
	}
	return cameras, nil
}
