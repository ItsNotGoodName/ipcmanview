package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/common"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

func NewDefaultEventHooks(bus *common.Bus, db repo.DB) DefaultEventHooks {
	return DefaultEventHooks{
		bus: bus,
		db:  db,
	}
}

type DefaultEventHooks struct {
	bus *common.Bus
	db  repo.DB
}

func (e DefaultEventHooks) logError(err error) {
	if err != nil {
		log.Err(err).Str("service", "dahua.DefaultEventHooks").Send()
	}
}

func (e DefaultEventHooks) Connecting(ctx context.Context, deviceID int64) {
	e.logError(e.db.CreateDahuaEventWorkerState(ctx, repo.CreateDahuaEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateConnecting,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerConnecting(models.EventDahuaEventWorkerConnecting{
		DeviceID: deviceID,
	})
}

func (e DefaultEventHooks) Connect(ctx context.Context, deviceID int64) {
	e.logError(e.db.CreateDahuaEventWorkerState(ctx, repo.CreateDahuaEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateConnected,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerConnect(models.EventDahuaEventWorkerConnect{
		DeviceID: deviceID,
	})
}

func (e DefaultEventHooks) Disconnect(ctx context.Context, deviceID int64, err error) {
	e.logError(e.db.CreateDahuaEventWorkerState(ctx, repo.CreateDahuaEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateDisconnected,
		Error:     repo.ErrorToNullString(err),
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerDisconnect(models.EventDahuaEventWorkerDisconnect{
		DeviceID: deviceID,
		Error:    err,
	})
}

func (e DefaultEventHooks) Event(ctx context.Context, event models.DahuaEvent) {
	eventRule, err := e.db.GetDahuaEventRuleByEvent(ctx, event)
	if err != nil {
		e.logError(err)
		return
	}

	deviceName, err := e.db.GetDahuaDeviceName(ctx, event.DeviceID)
	if err != nil && !repo.IsNotFound(err) {
		e.logError(err)
		return
	}

	if !eventRule.IgnoreDB {
		id, err := e.db.CreateDahuaEvent(ctx, repo.CreateDahuaEventParams{
			DeviceID:  event.DeviceID,
			Code:      event.Code,
			Action:    event.Action,
			Index:     int64(event.Index),
			Data:      event.Data,
			CreatedAt: types.NewTime(event.CreatedAt),
		})
		if err != nil {
			e.logError(err)
			return
		}
		event.ID = id
	}

	e.bus.EventDahuaEvent(models.EventDahuaEvent{
		DeviceName: deviceName,
		Event:      event,
		EventRule:  eventRule,
	})
}
