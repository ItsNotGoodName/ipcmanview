package dahuaweb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

func getEventRule(db repo.DB, evt models.DahuaEvent) (models.DahuaEventRule, error) {
	return models.DahuaEventRule{}, nil
}

// ---------- EventHooksProxy

func NewEventHooksProxy(bus *dahua.Bus, db repo.DB) EventHooksProxy {
	return EventHooksProxy{
		bus: bus,
		db:  db,
	}
}

// EventHooksProxy saves events into database.
type EventHooksProxy struct {
	bus *dahua.Bus
	db  repo.DB
}

func (p EventHooksProxy) CameraEvent(ctx context.Context, event models.DahuaEvent) {
	eventRule, err := getEventRule(p.db, event)
	if err != nil {
		log.Err(err).Msg("Failed to get DahuaEventRule")
		return
	}

	if !eventRule.IgnoreDB {
		id, err := p.db.CreateDahuaEvent(ctx, repo.CreateDahuaEventParams{
			CameraID:  event.CameraID,
			Code:      event.Code,
			Action:    event.Action,
			Index:     int64(event.Index),
			Data:      event.Data,
			CreatedAt: types.NewTime(event.CreatedAt),
		})
		if err != nil {
			log.Err(err).Msg("Failed to save DahuaEvent")
			return
		}
		event.ID = id
	}

	p.bus.CameraEvent(ctx, event, eventRule)
}

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

	return dbCamera.Convert(), true, nil
}

func (r Repo) ListConn(ctx context.Context) ([]models.DahuaConn, error) {
	dbCameras, err := r.db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	cameras := make([]models.DahuaConn, 0, len(dbCameras))
	for _, row := range dbCameras {
		cameras = append(cameras, row.Convert())
	}

	return cameras, nil
}
