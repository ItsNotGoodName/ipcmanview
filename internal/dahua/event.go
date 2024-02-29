package dahua

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
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

func publishEvent(ctx context.Context, db sqlite.DB, bus *event.Bus, deviceID int64, evt dahuacgi.Event) error {
	eventRule, err := getEventRuleByEvent(ctx, db, deviceID, evt.Code)
	if err != nil {
		return err
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
		id, err := db.C().DahuaCreateEvent(ctx, repo.DahuaCreateEventParams{
			DeviceID:  v.DeviceID,
			Code:      v.Code,
			Action:    v.Action,
			Index:     v.Index,
			Data:      v.Data,
			CreatedAt: v.CreatedAt,
		})
		if err != nil {
			return err
		}
		v.ID = id
	}

	bus.DahuaEvent(event.DahuaEvent{
		Event:     v,
		EventRule: eventRule,
	})

	return nil
}
