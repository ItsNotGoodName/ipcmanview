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

// ---------- EventHooksProxy

func NewEventHooksProxy(hooks dahua.EventHooks, db repo.DB) EventHooksProxy {
	return EventHooksProxy{
		hooks: hooks,
		db:    db,
	}
}

// EventHooksProxy saves events into database.
type EventHooksProxy struct {
	hooks dahua.EventHooks
	db    repo.DB
}

func (p EventHooksProxy) CameraEvent(ctx context.Context, evt models.DahuaEvent) {
	id, err := p.db.CreateDahuaEvent(ctx, repo.CreateDahuaEventParams{
		CameraID:  evt.CameraID,
		Code:      evt.Code,
		Action:    evt.Action,
		Index:     int64(evt.Index),
		Data:      evt.Data,
		CreatedAt: types.NewTime(evt.CreatedAt),
	})
	if err != nil {
		log.Err(err).Msg("Failed to save DahuaEvent")
		return
	}
	evt.ID = id
	p.hooks.CameraEvent(ctx, evt)
}

// ---------- CameraStore

func NewCameraStore(db repo.DB) CameraStore {
	return CameraStore{
		db: db,
	}
}

type CameraStore struct {
	db repo.DB
}

func (s CameraStore) Save(ctx context.Context, camera ...models.DahuaCamera) error {
	return errors.ErrUnsupported
}

func (s CameraStore) Get(ctx context.Context, id int64) (models.DahuaCamera, bool, error) {
	dbCamera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DahuaCamera{}, false, nil
		}
		return models.DahuaCamera{}, false, err
	}

	return dbCamera.Convert(), true, nil
}

func (s CameraStore) List(ctx context.Context) ([]models.DahuaCamera, error) {
	dbCameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	cameras := make([]models.DahuaCamera, 0, len(dbCameras))
	for _, row := range dbCameras {
		cameras = append(cameras, row.Convert())
	}

	return cameras, nil
}
