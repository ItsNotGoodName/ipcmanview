package core

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func NewLocation(location string) (models.Location, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return models.Location{}, err
	}

	return models.Location{
		Location: loc,
	}, nil
}

type DahuaCameraStore interface {
	List(ctx context.Context) ([]models.DahuaCamera, error)
	Save(ctx context.Context, camera ...models.DahuaCamera) error
}

func DahuaBootstrap(ctx context.Context, cameraStore DahuaCameraStore, store *dahua.Store, eventWorkerStore *dahua.EventWorkerStore) error {
	cameras, err := cameraStore.List(ctx)
	if err != nil {
		return err
	}
	conns := store.ConnList(ctx, cameras)
	for _, conn := range conns {
		if err := eventWorkerStore.Create(conn.Camera); err != nil {
			return err
		}
	}
	return err
}
