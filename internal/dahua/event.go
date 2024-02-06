package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/rs/zerolog/log"
)

func NewDefaultEventHooks(bus *event.Bus, db sqlite.DB) DefaultEventHooks {
	return DefaultEventHooks{
		bus: bus,
		db:  db,
	}
}

type DefaultEventHooks struct {
	bus *event.Bus
	db  sqlite.DB
}

func (e DefaultEventHooks) logError(err error) {
	if err != nil {
		log.Err(err).Str("service", "dahua.DefaultEventHooks").Send()
	}
}

func (e DefaultEventHooks) Connecting(ctx context.Context, deviceID int64) {
	e.logError(e.db.C().DahuaCreateEventWorkerState(ctx, repo.DahuaCreateEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateConnecting,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.DahuaEventWorkerConnecting(event.DahuaEventWorkerConnecting{
		DeviceID: deviceID,
	})
}

func (e DefaultEventHooks) Connect(ctx context.Context, deviceID int64) {
	e.logError(e.db.C().DahuaCreateEventWorkerState(ctx, repo.DahuaCreateEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateConnected,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.DahuaEventWorkerConnect(event.DahuaEventWorkerConnect{
		DeviceID: deviceID,
	})
}

func (e DefaultEventHooks) Disconnect(ctx context.Context, deviceID int64, err error) {
	e.logError(e.db.C().DahuaCreateEventWorkerState(ctx, repo.DahuaCreateEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateDisconnected,
		Error:     core.ErrorToNullString(err),
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.DahuaEventWorkerDisconnect(event.DahuaEventWorkerDisconnect{
		DeviceID: deviceID,
		Error:    err,
	})
}

func (e DefaultEventHooks) Event(ctx context.Context, deviceID int64, evt dahuacgi.Event) {
	eventRule, err := GetEventRuleByEvent(ctx, e.db, deviceID, evt.Code)
	if err != nil {
		e.logError(err)
		return
	}

	v := repo.DahuaEvent{
		ID:        0,
		DeviceID:  deviceID,
		Code:      evt.Code,
		Action:    evt.Action,
		Index:     int64(evt.Index),
		Data:      types.NewJSON(evt.Data),
		CreatedAt: types.NewTime(time.Now()),
	}
	if !eventRule.IgnoreDb {
		id, err := e.db.C().DahuaCreateEvent(ctx, repo.DahuaCreateEventParams{
			DeviceID:  v.DeviceID,
			Code:      v.Code,
			Action:    v.Action,
			Index:     v.Index,
			Data:      v.Data,
			CreatedAt: v.CreatedAt,
		})
		if err != nil {
			e.logError(err)
			return
		}
		v.ID = id
	}

	e.bus.DahuaEvent(event.DahuaEvent{
		Event:     v,
		EventRule: eventRule,
	})
}
