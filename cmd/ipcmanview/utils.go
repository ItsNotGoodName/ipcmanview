package main

import (
	"crypto/rand"
	"os"
	"path/filepath"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/server"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/spf13/afero"
)

func mkdir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

type Shared struct {
	Dir string `default:"ipcmanview_data" env:"DIR" help:"Directory path for storing data."`
}

func (c *Shared) init() error {
	var err error
	c.Dir, err = filepath.Abs(c.Dir)
	if err != nil {
		return err
	}
	return os.MkdirAll(c.Dir, 0755)
}

func (c Shared) useDahuaAFS() (afero.Fs, error) {
	dir := filepath.Join(c.Dir, "dahua-files")
	if err := mkdir(dir); err != nil {
		return nil, err
	}
	return afero.NewBasePathFs(afero.NewOsFs(), dir), nil
}

func (c Shared) useSecret() ([]byte, error) {
	dir := filepath.Join(c.Dir, "secret")

	b, err := os.ReadFile(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		b = make([]byte, 64)

		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(dir, b, 0600); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (c Shared) useDB(ctx *Context) (repo.DB, error) {
	sqlDB, err := sqlite.New(filepath.Join(c.Dir, "sqlite.db"))
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
		return repo.DB{}, err
	}

	return db, nil
}

func (c Shared) useCert() (server.Certificate, error) {
	cert := server.Certificate{
		CertFile: filepath.Join(c.Dir, "cert.pem"),
		KeyFile:  filepath.Join(c.Dir, "key.pem"),
	}

	certFileExists, err := core.FileExists(cert.CertFile)
	if err != nil {
		return cert, err
	}
	keyFileExists, err := core.FileExists(cert.KeyFile)
	if err != nil {
		return cert, err
	}
	if !(certFileExists && keyFileExists) {
		err := server.GenerateCertificate(cert)
		if err != nil {
			return cert, err
		}
	}

	return cert, nil
}

// type SharedDevices struct {
// 	ID  []int64 `help:"Run on device by ID."`
// 	All bool    `help:"Run on all devices."`
// }

// func (c SharedDevices) useDevices(ctx context.Context, db repo.DB) ([]models.DahuaDeviceConn, error) {
// 	var devices []models.DahuaDeviceConn
// 	if c.All {
// 		dbDevices, err := db.DahuaListDevices(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		for _, dbDevice := range dbDevices {
// 			devices = append(devices, dbDevice.Convert())
// 		}
// 	} else {
// 		dbDevices, err := db.ListDahuaDevicesByIDs(ctx, c.ID)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		for _, dbDevice := range dbDevices {
// 			devices = append(devices, dbDevice.Convert())
// 		}
// 	}
// 	return devices, nil
// }
