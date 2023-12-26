package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func NewDefaultEventHooks(bus *core.Bus, db repo.DB) DefaultEventHooks {
	return DefaultEventHooks{
		ServiceContext: sutureext.NewServiceContext("dahua.DefailtEventHooks"),
		bus:            bus,
		db:             db,
	}
}

type DefaultEventHooks struct {
	sutureext.ServiceContext
	bus *core.Bus
	db  repo.DB
}

func (e DefaultEventHooks) logErr(err error) {
	if err != nil {
		log.Err(err).Str("service", e.String()).Send()
	}
}

func (e DefaultEventHooks) Connecting(ctx context.Context, deviceID int64) {
	e.logErr(e.db.CreateDahuaEventWorkerState(e.Context(), repo.CreateDahuaEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateConnecting,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerConnecting(models.EventDahuaEventWorkerConnecting{
		DeviceID: deviceID,
	})
}

func (e DefaultEventHooks) Connect(ctx context.Context, deviceID int64) {
	e.logErr(e.db.CreateDahuaEventWorkerState(e.Context(), repo.CreateDahuaEventWorkerStateParams{
		DeviceID:  deviceID,
		State:     models.DahuaEventWorkerStateConnected,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerConnect(models.EventDahuaEventWorkerConnect{
		DeviceID: deviceID,
	})
}

func (e DefaultEventHooks) Disconnect(deviceID int64, err error) {
	e.logErr(e.db.CreateDahuaEventWorkerState(e.Context(), repo.CreateDahuaEventWorkerStateParams{
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
		e.logErr(err)
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
			e.logErr(err)
			return
		}
		event.ID = id
	}

	e.bus.EventDahuaDeviceEvent(models.EventDahuaDeviceEvent{
		Event:     event,
		EventRule: eventRule,
	})
}
