package webdahua

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/rs/zerolog/log"
)

// ---------- Convert

func ConvertListDahuaCameraRows(dbCameras []sqlc.ListDahuaCameraRow) []models.DahuaCamera {
	cameras := make([]models.DahuaCamera, 0, len(dbCameras))
	for _, c := range dbCameras {
		cameras = append(cameras, models.DahuaCamera{
			ID:        c.ID,
			Address:   c.Address,
			Username:  c.Username,
			Password:  c.Password,
			Location:  c.Location,
			Seed:      int(c.Seed),
			CreatedAt: c.CreatedAt,
		})
	}
	return cameras
}

func ConvertGetDahuaCameraRow(c sqlc.GetDahuaCameraRow) models.DahuaCamera {
	return models.DahuaCamera{
		ID:        c.ID,
		Address:   c.Address,
		Username:  c.Username,
		Password:  c.Password,
		Location:  c.Location,
		Seed:      int(c.Seed),
		CreatedAt: c.CreatedAt,
	}
}

// ---------- DahuaEventHooksProxy

func NewDahuaEventHooksProxy(hooks dahua.EventHooks, db sqlc.DB) DahuaEventHooksProxy {
	return DahuaEventHooksProxy{
		hooks: hooks,
		db:    db,
	}
}

type DahuaEventHooksProxy struct {
	hooks dahua.EventHooks
	db    sqlc.DB
}

func (p DahuaEventHooksProxy) CameraEvent(evt models.DahuaEvent) {
	id, err := p.db.CreateDahuaEvent(context.TODO(), sqlc.CreateDahuaEventParams{
		CameraID:      evt.CameraID,
		ContentType:   evt.ContentType,
		ContentLength: int64(evt.ContentLength),
		Code:          evt.Code,
		Action:        evt.Action,
		Index:         int64(evt.Index),
		Data:          evt.Data,
		CreatedAt:     evt.CreatedAt,
	})
	if err != nil {
		log.Err(err).Caller().Msg("Failed to save DahuaEvent")
		return
	}
	evt.ID = id
	p.hooks.CameraEvent(evt)
}

// ---------- DahuaCameraStore

func NewDahuaCameraStore(db sqlc.DB) DahuaCameraStore {
	return DahuaCameraStore{
		db: db,
	}
}

type DahuaCameraStore struct {
	db sqlc.DB
}

func (s DahuaCameraStore) Save(ctx context.Context, camera ...models.DahuaCamera) error {
	for _, camera := range camera {
		now := time.Now()
		s.db.UpsertDahuaCamera(ctx, camera.ID, sqlc.CreateDahuaCameraParams{
			Name:      camera.Address,
			Address:   camera.Address,
			Username:  camera.Username,
			Password:  camera.Password,
			Location:  camera.Location,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return nil
}

func (s DahuaCameraStore) Get(ctx context.Context, id int64) (models.DahuaCamera, bool, error) {
	camera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DahuaCamera{}, false, nil
		}
		return models.DahuaCamera{}, false, err
	}

	return ConvertGetDahuaCameraRow(camera), true, nil
}

func (s DahuaCameraStore) List(ctx context.Context) ([]models.DahuaCamera, error) {
	dbCameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	return ConvertListDahuaCameraRows(dbCameras), nil
}
