package dahua

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

// ---------- Repo

func NewRepo(db repo.DB) Repo {
	return Repo{
		db: db,
	}
}

type Repo struct {
	db repo.DB
}

func (r Repo) GetFileByFilePath(ctx context.Context, cameraID int64, filePath string) (models.DahuaFile, error) {
	file, err := r.db.GetDahuaFileByFilePath(ctx, repo.GetDahuaFileByFilePathParams{
		CameraID: cameraID,
		FilePath: filePath,
	})
	if err != nil {
		return models.DahuaFile{}, err
	}

	return file.Convert(), nil
}

func (r Repo) GetConn(ctx context.Context, id int64) (models.DahuaConn, bool, error) {
	dbCamera, err := r.db.GetDahuaCamera(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DahuaConn{}, false, nil
		}
		return models.DahuaConn{}, false, err
	}

	return dbCamera.Convert().DahuaConn, true, nil
}

func (r Repo) ListConn(ctx context.Context) ([]models.DahuaConn, error) {
	dbCameras, err := r.db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	cameras := make([]models.DahuaConn, 0, len(dbCameras))
	for _, row := range dbCameras {
		cameras = append(cameras, row.Convert().DahuaConn)
	}

	return cameras, nil
}
