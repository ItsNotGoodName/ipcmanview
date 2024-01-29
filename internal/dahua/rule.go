package dahua

import (
	"context"
	"errors"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateEventRule(ctx context.Context, db repo.DB, arg repo.DahuaCreateEventRuleParams) error {
	arg.Code = strings.TrimSpace(arg.Code)
	if arg.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.DahuaCreateEventRule(ctx, arg)
}

func UpdateEventRule(ctx context.Context, db repo.DB, rule repo.DahuaEventRule, arg repo.DahuaUpdateEventRuleParams) error {
	arg.Code = strings.TrimSpace(arg.Code)
	if rule.Code == "" {
		arg.Code = rule.Code
	}

	if arg.Code == "" && rule.Code != "" {
		return errors.New("code cannot be empty")
	}

	return db.DahuaUpdateEventRule(ctx, arg)
}

func DeleteEventRule(ctx context.Context, db repo.DB, rule repo.DahuaEventRule) error {
	if rule.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.DahuaDeleteEventRule(ctx, rule.ID)
}

func GetEventRuleByEvent(ctx context.Context, db repo.DB, deviceID int64, code string) (repo.DahuaEventRule, error) {
	res, err := db.DahuaGetEventRuleByEvent(ctx, repo.DahuaGetEventRuleByEventParams{
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
