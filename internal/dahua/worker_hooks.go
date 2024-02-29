package dahua

import (
	"context"
	"database/sql"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

func NewDefaultWorkerHooks(db sqlite.DB, bus *event.Bus) DefaultWorkerHooks {
	return DefaultWorkerHooks{
		db:  db,
		bus: bus,
	}
}

type DefaultWorkerHooks struct {
	db  sqlite.DB
	bus *event.Bus
}

func (h DefaultWorkerHooks) Serve(ctx context.Context, w Worker, connected bool, fn func(ctx context.Context) error) error {
	state := models.DahuaWorkerStateConnecting
	if connected {
		state = models.DahuaWorkerStateConnected
	}
	err := h.db.C().DahuaCreateWorkerEvent(ctx, repo.DahuaCreateWorkerEventParams{
		DeviceID:  w.DeviceID,
		Type:      w.Type,
		State:     state,
		Error:     sql.NullString{},
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}

	if connected {
		h.bus.DahuaWorkerConnected(event.DahuaWorkerConnected{
			DeviceID: w.DeviceID,
			Type:     w.Type,
		})
	} else {
		h.bus.DahuaWorkerConnecting(event.DahuaWorkerConnecting{
			DeviceID: w.DeviceID,
			Type:     w.Type,
		})
	}

	serveError := fn(ctx)

	err = h.db.C().DahuaCreateWorkerEvent(ctx, repo.DahuaCreateWorkerEventParams{
		DeviceID:  w.DeviceID,
		Type:      w.Type,
		State:     models.DahuaWorkerStateDisconnected,
		Error:     core.ErrorToNullString(serveError),
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}
	h.bus.DahuaWorkerDisconnected(event.DahuaWorkerDisconnected{
		DeviceID: w.DeviceID,
		Type:     w.Type,
	})

	return serveError
}

func (h DefaultWorkerHooks) Connected(ctx context.Context, w Worker) {
	err := h.db.C().DahuaCreateWorkerEvent(ctx, repo.DahuaCreateWorkerEventParams{
		DeviceID:  w.DeviceID,
		Type:      w.Type,
		State:     models.DahuaWorkerStateConnected,
		Error:     sql.NullString{},
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		log.Err(err).Send()
		return
	}
	h.bus.DahuaWorkerConnected(event.DahuaWorkerConnected{
		DeviceID: w.DeviceID,
		Type:     w.Type,
	})
}
