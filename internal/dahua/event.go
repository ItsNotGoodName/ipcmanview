package dahua

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
)

const eventRuleCodeErrorMessage = "Code cannot be empty."

func CreateEventRule(ctx context.Context, arg repo.DahuaCreateEventRuleParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	// Mutate
	arg.Code = strings.TrimSpace(arg.Code)

	if arg.Code == "" {
		return 0, core.NewFieldError("Code", eventRuleCodeErrorMessage)
	}

	return app.DB.C().DahuaCreateEventRule(ctx, arg)
}

func UpdateEventRule(ctx context.Context, arg repo.DahuaUpdateEventRuleParams) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	// Mutate
	arg.Code = strings.TrimSpace(arg.Code)

	model, err := app.DB.C().DahuaGetEventRule(ctx, arg.ID)
	if err != nil {
		return err
	}

	if model.Code == "" && arg.Code != "" {
		return core.NewFieldError("Code", eventRuleCodeErrorMessage)
	}

	return app.DB.C().DahuaUpdateEventRule(ctx, arg)
}

func DeleteEventRule(ctx context.Context, id int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	model, err := app.DB.C().DahuaGetEventRule(ctx, id)
	if err != nil {
		return err
	}
	if model.Code == "" {
		return core.ErrForbidden
	}

	return app.DB.C().DahuaDeleteEventRule(ctx, model.ID)
}

func getEventRuleByEvent(ctx context.Context, deviceID int64, code string) (repo.DahuaEventRule, error) {
	res, err := app.DB.C().DahuaGetEventRuleByEvent(ctx, repo.DahuaGetEventRuleByEventParams{
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

func publishEvent(ctx context.Context, deviceID int64, event dahuacgi.Event) error {
	eventRule, err := getEventRuleByEvent(ctx, deviceID, event.Code)
	if err != nil {
		return err
	}

	v := repo.DahuaEvent{
		ID:        0,
		DeviceID:  deviceID,
		Code:      event.Code,
		Action:    event.Action,
		Index:     int64(event.Index),
		Data:      types.NewJSON(core.IgnoreError(json.MarshalIndent(event.Data, "", "  "))),
		CreatedAt: types.NewTime(time.Now()),
	}
	if !eventRule.IgnoreDb {
		id, err := app.DB.C().DahuaCreateEvent(ctx, repo.DahuaCreateEventParams{
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

	app.Hub.DahuaEvent(bus.DahuaEvent{
		Event:     v,
		EventRule: eventRule,
	})

	return nil
}
