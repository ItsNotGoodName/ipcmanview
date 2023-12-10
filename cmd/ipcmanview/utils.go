package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

type Shared struct {
	DBPath string `default:"sqlite.db" env:"DB_PATH" help:"Path to SQLite database."`
}

func useDB(ctx *Context, path string) (sqlc.DB, error) {
	sqlDB, err := sqlite.New(path)
	if err != nil {
		return sqlc.DB{}, err
	}
	if err := migrations.Migrate(sqlDB); err != nil {
		return sqlc.DB{}, err
	}
	if ctx.Debug {
		return sqlc.NewDB(sqlite.NewDebugDB(sqlDB)), nil
	}
	return sqlc.NewDB(sqlite.NewDB(sqlDB)), nil
}

type SharedCameras struct {
	ID  []int64 `help:"Run on camera by ID."`
	All bool    `help:"Run on all cameras."`
}

func (c SharedCameras) useCameras(ctx context.Context, db sqlc.DB) ([]models.DahuaCameraInfo, error) {
	var cameras []models.DahuaCameraInfo
	if c.All {
		dbCameras, err := db.ListDahuaCamera(ctx)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range dbCameras {
			cameras = append(cameras, models.DahuaCameraInfo{
				DahuaCamera: dbCamera.Convert(),
				Name:        dbCamera.Name,
				UpdatedAt:   dbCamera.UpdatedAt.Time,
			})
		}
	} else {
		dbCameras, err := db.ListDahuaCameraByIDs(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range dbCameras {
			cameras = append(cameras, models.DahuaCameraInfo{
				DahuaCamera: dbCamera.Convert(),
				Name:        dbCamera.Name,
				UpdatedAt:   dbCamera.UpdatedAt.Time,
			})
		}
	}
	return cameras, nil
}
