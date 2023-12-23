package dahua

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func NewRepo(db repo.DB) Repo {
	return Repo{
		db: db,
	}
}

type Repo struct {
	db repo.DB
}

func (r Repo) GetFileByFilePath(ctx context.Context, deviceID int64, filePath string) (models.DahuaFile, error) {
	file, err := r.db.GetDahuaFileByFilePath(ctx, repo.GetDahuaFileByFilePathParams{
		DeviceID: deviceID,
		FilePath: filePath,
	})
	if err != nil {
		return models.DahuaFile{}, err
	}

	return file.Convert(), nil
}

func (r Repo) GetConn(ctx context.Context, id int64) (models.DahuaConn, bool, error) {
	dbDevice, err := r.db.GetDahuaDevice(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DahuaConn{}, false, nil
		}
		return models.DahuaConn{}, false, err
	}

	return dbDevice.Convert().DahuaConn, true, nil
}

func (r Repo) ListConn(ctx context.Context) ([]models.DahuaConn, error) {
	dbDevices, err := r.db.ListDahuaDevice(ctx)
	if err != nil {
		return nil, err
	}

	devices := make([]models.DahuaConn, 0, len(dbDevices))
	for _, row := range dbDevices {
		devices = append(devices, row.Convert().DahuaConn)
	}

	return devices, nil
}
