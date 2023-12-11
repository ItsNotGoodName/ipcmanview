package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

type SharedDB struct {
	DBPath string `default:"sqlite.db" env:"DB_PATH" help:"Path to SQLite database."`
}

func (c SharedDB) useDB(ctx *Context) (repo.DB, error) {
	sqlDB, err := sqlite.New(c.DBPath)
	if err != nil {
		return repo.DB{}, err
	}
	if err := migrations.Migrate(sqlDB); err != nil {
		return repo.DB{}, err
	}
	if ctx.Debug {
		return repo.NewDB(sqlite.NewDebugDB(sqlDB)), nil
	}
	return repo.NewDB(sqlite.NewDB(sqlDB)), nil
}

type SharedCameras struct {
	ID  []int64 `help:"Run on camera by ID."`
	All bool    `help:"Run on all cameras."`
}

func (c SharedCameras) useCameras(ctx context.Context, db repo.DB) ([]models.DahuaCameraConn, error) {
	var cameras []models.DahuaCameraConn
	if c.All {
		dbCameras, err := db.ListDahuaCamera(ctx)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range dbCameras {
			cameras = append(cameras, models.DahuaCameraConn{
				DahuaConn: dbCamera.Convert(),
			})
		}
	} else {
		dbCameras, err := db.ListDahuaCameraByIDs(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range dbCameras {
			cameras = append(cameras, models.DahuaCameraConn{
				DahuaConn: dbCamera.Convert(),
			})
		}
	}
	return cameras, nil
}
