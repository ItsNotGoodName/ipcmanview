package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/spf13/afero"
)

type Shared struct {
	Dir string `default:"ipcmanview_data" env:"DIR" help:"Directory path for storing data."`
}

func mkdir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func (c Shared) useDir() (string, error) {
	dir, err := filepath.Abs(c.Dir)
	if err != nil {
		return "", err
	}

	if err := mkdir(dir); err != nil {
		return "", err
	}

	return dir, nil
}

func (c Shared) useDahuaAFS() (afero.Fs, error) {
	dir, err := c.useDir()
	if err != nil {
		return nil, err
	}
	dir = filepath.Join(dir, "dahua-files")
	if err := mkdir(dir); err != nil {
		return nil, err
	}

	return afero.NewBasePathFs(afero.NewOsFs(), dir), nil
}

func (c Shared) useDB(ctx *Context) (repo.DB, error) {
	dir, err := c.useDir()
	if err != nil {
		return repo.DB{}, err
	}

	sqlDB, err := sqlite.New(filepath.Join(dir, "sqlite.db"))
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

	err = migrations.Normalize(ctx, db)
	if err != nil {
		return repo.DB{}, nil
	}

	return db, nil
}

type SharedDevices struct {
	ID  []int64 `help:"Run on device by ID."`
	All bool    `help:"Run on all devices."`
}

func (c SharedDevices) useDevices(ctx context.Context, db repo.DB) ([]models.DahuaDeviceConn, error) {
	var devices []models.DahuaDeviceConn
	if c.All {
		dbDevices, err := db.ListDahuaDevice(ctx)
		if err != nil {
			return nil, err
		}

		for _, dbDevice := range dbDevices {
			devices = append(devices, dbDevice.Convert())
		}
	} else {
		dbDevices, err := db.ListDahuaDeviceByIDs(ctx, c.ID)
		if err != nil {
			return nil, err
		}

		for _, dbDevice := range dbDevices {
			devices = append(devices, dbDevice.Convert())
		}
	}
	return devices, nil
}
