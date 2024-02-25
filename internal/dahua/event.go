package dahua

import (
	"context"
	"encoding/json"
	"strings"
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

const eventRuleCodeErrorMessage = "Code cannot be empty."

func CreateEventRule(ctx context.Context, db sqlite.DB, arg repo.DahuaCreateEventRuleParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	// Mutate
	arg.Code = strings.TrimSpace(arg.Code)

	if arg.Code == "" {
		return 0, core.NewFieldError("Code", eventRuleCodeErrorMessage)
	}

	return db.C().DahuaCreateEventRule(ctx, arg)
}

func UpdateEventRule(ctx context.Context, db sqlite.DB, arg repo.DahuaUpdateEventRuleParams) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	// Mutate
	arg.Code = strings.TrimSpace(arg.Code)

	model, err := db.C().DahuaGetEventRule(ctx, arg.ID)
	if err != nil {
		return err
	}

	if model.Code == "" && arg.Code != "" {
		return core.NewFieldError("Code", eventRuleCodeErrorMessage)
	}

	return db.C().DahuaUpdateEventRule(ctx, arg)
}

func DeleteEventRule(ctx context.Context, db sqlite.DB, id int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	model, err := db.C().DahuaGetEventRule(ctx, id)
	if err != nil {
		return err
	}
	if model.Code == "" {
		return core.ErrForbidden
	}

	return db.C().DahuaDeleteEventRule(ctx, model.ID)
}

func getEventRuleByEvent(ctx context.Context, db sqlite.DB, deviceID int64, code string) (repo.DahuaEventRule, error) {
	res, err := db.C().DahuaGetEventRuleByEvent(ctx, repo.DahuaGetEventRuleByEventParams{
		DeviceID: deviceID,
		Code:     code,
	})
	if err != nil {
		return repo.DahuaEventRule{}, err
	}
	if len(res) == 0 {
		return repo.DahuaEventRule{}, nil
	}

	return repo.DahuaEventRule{
		ID:         0,
		Code:       code,
		IgnoreDb:   res[0].IgnoreDb,
		IgnoreLive: res[0].IgnoreLive,
		IgnoreMqtt: res[0].IgnoreMqtt,
	}, nil
}

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
	eventRule, err := getEventRuleByEvent(ctx, e.db, deviceID, evt.Code)
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
		Data:      types.NewJSON(core.IgnoreError(json.MarshalIndent(evt.Data, "", "  "))),
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
