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

func NewEventHooks(bus *core.Bus, db repo.DB) EventHooks {
	return EventHooks{
		ServiceContext: sutureext.NewServiceContext("dahua.EventHooks"),
		bus:            bus,
		db:             db,
	}
}

type EventHooks struct {
	sutureext.ServiceContext
	bus *core.Bus
	db  repo.DB
}

func (e EventHooks) logErr(err error) {
	if err != nil {
		log.Err(err).Str("service", e.String()).Send()
	}
}

func (e EventHooks) Connecting(ctx context.Context, cameraID int64) {
	e.logErr(e.db.CreateDahuaEventWorkerState(e.Context(), repo.CreateDahuaEventWorkerStateParams{
		CameraID:  cameraID,
		State:     models.DahuaEventWorkerStateConnecting,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerConnecting(models.EventDahuaEventWorkerConnecting{
		CameraID: cameraID,
	})
}

func (e EventHooks) Connect(ctx context.Context, cameraID int64) {
	e.logErr(e.db.CreateDahuaEventWorkerState(e.Context(), repo.CreateDahuaEventWorkerStateParams{
		CameraID:  cameraID,
		State:     models.DahuaEventWorkerStateConnected,
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerConnect(models.EventDahuaEventWorkerConnect{
		CameraID: cameraID,
	})
}

func (e EventHooks) Disconnect(cameraID int64, err error) {
	e.logErr(e.db.CreateDahuaEventWorkerState(e.Context(), repo.CreateDahuaEventWorkerStateParams{
		CameraID:  cameraID,
		State:     models.DahuaEventWorkerStateDisconnected,
		Error:     repo.ErrorToNullString(err),
		CreatedAt: types.NewTime(time.Now()),
	}))
	e.bus.EventDahuaEventWorkerDisconnect(models.EventDahuaEventWorkerDisconnect{
		CameraID: cameraID,
		Error:    err,
	})
}

func (e EventHooks) Event(ctx context.Context, event models.DahuaEvent) {
	eventRule, err := e.db.GetDahuaEventRuleByEvent(ctx, event)
	if err != nil {
		e.logErr(err)
		return
	}

	if !eventRule.IgnoreDB {
		id, err := e.db.CreateDahuaEvent(ctx, repo.CreateDahuaEventParams{
			CameraID:  event.CameraID,
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

	e.bus.EventDahuaCameraEvent(models.EventDahuaCameraEvent{
		Event:     event,
		EventRule: eventRule,
	})
}
