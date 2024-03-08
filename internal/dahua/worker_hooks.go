package dahua

import (
	"context"
	"database/sql"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

func NewDefaultWorkerHooks() DefaultWorkerHooks {
	return DefaultWorkerHooks{}
}

type DefaultWorkerHooks struct {
}

func (h DefaultWorkerHooks) Serve(ctx context.Context, w Worker, connected bool, fn func(ctx context.Context) error) error {
	state := models.DahuaWorkerState_Connecting
	if connected {
		state = models.DahuaWorkerState_Connected
	}
	err := app.DB.C().DahuaCreateWorkerEvent(ctx, repo.DahuaCreateWorkerEventParams{
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
		app.Hub.DahuaWorkerConnected(bus.DahuaWorkerConnected{
			DeviceID: w.DeviceID,
			Type:     w.Type,
		})
	} else {
		app.Hub.DahuaWorkerConnecting(bus.DahuaWorkerConnecting{
			DeviceID: w.DeviceID,
			Type:     w.Type,
		})
	}

	serveError := fn(ctx)

	err = app.DB.C().DahuaCreateWorkerEvent(ctx, repo.DahuaCreateWorkerEventParams{
		DeviceID:  w.DeviceID,
		Type:      w.Type,
		State:     models.DahuaWorkerState_Disconnected,
		Error:     core.ErrorToNullString(serveError),
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}
	app.Hub.DahuaWorkerDisconnected(bus.DahuaWorkerDisconnected{
		DeviceID: w.DeviceID,
		Type:     w.Type,
		Error:    err,
	})

	return serveError
}

func (h DefaultWorkerHooks) Connected(ctx context.Context, w Worker) {
	err := app.DB.C().DahuaCreateWorkerEvent(ctx, repo.DahuaCreateWorkerEventParams{
		DeviceID:  w.DeviceID,
		Type:      w.Type,
		State:     models.DahuaWorkerState_Connected,
		Error:     sql.NullString{},
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		log.Err(err).Send()
		return
	}
	app.Hub.DahuaWorkerConnected(bus.DahuaWorkerConnected{
		DeviceID: w.DeviceID,
		Type:     w.Type,
	})
}
