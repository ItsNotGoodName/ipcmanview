package webdahua

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

var Locations []string

//go:embed locations.txt
var locationsStr string

func init() {

	for _, location := range strings.Split(locationsStr, "\n") {
		if location != "" {
			Locations = append(Locations, location)
		}
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

func (p DahuaEventHooksProxy) CameraEvent(ctx context.Context, evt models.DahuaEvent) {
	id, err := p.db.CreateDahuaEvent(ctx, sqlc.CreateDahuaEventParams{
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
	return errors.ErrUnsupported
}

func (s DahuaCameraStore) Get(ctx context.Context, id int64) (models.DahuaCamera, bool, error) {
	dbCamera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DahuaCamera{}, false, nil
		}
		return models.DahuaCamera{}, false, err
	}

	return dbCamera.Convert(), true, nil
}

func (s DahuaCameraStore) List(ctx context.Context) ([]models.DahuaCamera, error) {
	dbCameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	return sqlc.ConvertListDahuaCameraRow(dbCameras), nil
}
