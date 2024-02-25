package dahua

import (
	"context"
	"errors"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

func CreateEventRule(ctx context.Context, db sqlite.DB, arg repo.DahuaCreateEventRuleParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}
	arg.Code = strings.TrimSpace(arg.Code)
	if arg.Code == "" {
		return 0, core.NewFieldError("Code", "code cannot be empty")
	}
	return db.C().DahuaCreateEventRule(ctx, arg)
}

func UpdateEventRule(ctx context.Context, db sqlite.DB, arg repo.DahuaUpdateEventRuleParams) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	rule, err := db.C().DahuaGetEventRule(ctx, arg.ID)
	if err != nil {
		return err
	}

	arg.Code = strings.TrimSpace(arg.Code)
	if rule.Code == "" {
		arg.Code = rule.Code
	}
	if arg.Code == "" && rule.Code != "" {
		return core.NewFieldError("Code", "code cannot be empty")
	}

	return db.C().DahuaUpdateEventRule(ctx, arg)
}

func DeleteEventRule(ctx context.Context, db sqlite.DB, id int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	rule, err := db.C().DahuaGetEventRule(ctx, id)
	if err != nil {
		return err
	}
	if rule.Code == "" {
		return errors.New("code cannot be empty")
	}

	return db.C().DahuaDeleteEventRule(ctx, rule.ID)
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
