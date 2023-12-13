package main

import (
	"context"
	"os"
	"path"

	"github.com/ItsNotGoodName/ipcmanview/internal/files"
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

type Shared struct {
	Dir string `default:"ipcmanview_data" env:"DIR" help:"Directory for data."`
}

func useDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func (c Shared) useDahuaFileStore() (files.DahuaFileStore, error) {
	dir := path.Join(c.Dir, "dahua-files")
	return files.NewDahuaFileStore(dir), useDir(dir)
}

func (c Shared) useDB(ctx *Context) (repo.DB, error) {
	if err := useDir(c.Dir); err != nil {
		return repo.DB{}, err
	}

	sqlDB, err := sqlite.New(path.Join(c.Dir, "sqlite.db"))
	if err != nil {
		return repo.DB{}, err
	}
	if err := migrations.Migrate(sqlDB); err != nil {
		return repo.DB{}, err
	}

	var db repo.DB
	if ctx.Debug {
		db = repo.NewDB(sqlite.NewDebugDB(sqlDB))
	} else {
		db = repo.NewDB(sqlite.NewDB(sqlDB))
	}

	return db, nil
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
			cameras = append(cameras, dbCamera.Convert())
		}
	} else {
		dbCameras, err := db.ListDahuaCameraByIDs(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		for _, dbCamera := range dbCameras {
			cameras = append(cameras, dbCamera.Convert())
		}
	}
	return cameras, nil
}
