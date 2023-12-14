package dahua

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

func NewEventHooks(bus *dahuacore.Bus, db repo.DB) EventHooks {
	return EventHooks{
		bus: bus,
		db:  db,
	}
}

type EventHooks struct {
	bus *dahuacore.Bus
	db  repo.DB
}

func (EventHooks) logErr(err error) {
	if err != nil {
		log.Err(err).Str("package", "dahuacore").Msg("Failed to handle event")
	}
}

func (e EventHooks) Connecting(ctx context.Context, cameraID int64) {
	e.logErr(e.db.CreateDahuaEventWorkerState(context.Background(), repo.CreateDahuaEventWorkerStateParams{
		CameraID:  cameraID,
		State:     models.DahuaEventWorkerStateConnecting,
		CreatedAt: types.NewTime(time.Now()),
	}))
}

func (e EventHooks) Connected(ctx context.Context, cameraID int64) {
	e.logErr(e.db.CreateDahuaEventWorkerState(context.Background(), repo.CreateDahuaEventWorkerStateParams{
		CameraID:  cameraID,
		State:     models.DahuaEventWorkerStateConnected,
		CreatedAt: types.NewTime(time.Now()),
	}))
}

func (e EventHooks) Disconnect(cameraID int64, err error) {
	e.logErr(e.db.CreateDahuaEventWorkerState(context.Background(), repo.CreateDahuaEventWorkerStateParams{
		CameraID:  cameraID,
		State:     models.DahuaEventWorkerStateDisconnected,
		Error:     repo.ErrorToNullString(err),
		CreatedAt: types.NewTime(time.Now()),
	}))
}

func (e EventHooks) Event(ctx context.Context, event models.DahuaEvent) {
	eventRule, err := e.db.GetDahuaEventRuleByEvent(ctx, event)
	if err != nil {
		log.Err(err).Msg("Failed to get DahuaEventRule")
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
			log.Err(err).Msg("Failed to save DahuaEvent")
			return
		}
		event.ID = id
	}

	e.bus.CameraEvent(ctx, event, eventRule)
}
