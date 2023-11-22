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

func DahuaBootstrap(ctx context.Context, store *dahua.Store, eventWorkerStore *dahua.EventWorkerStore) error {
	conns, err := store.ConnList(ctx)
	if err != nil {
		return err
	}
	for _, conn := range conns {
		if err := eventWorkerStore.Create(conn.Camera); err != nil {
			return err
		}
	}
	return err
}
