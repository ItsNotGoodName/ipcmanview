package dahua

import (
	"context"
	"errors"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateEventRule(ctx context.Context, db repo.DB, arg repo.CreateDahuaEventRuleParams) error {
	arg.Code = strings.TrimSpace(arg.Code)
	if arg.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.CreateDahuaEventRule(ctx, arg)
}

func UpdateEventRule(ctx context.Context, db repo.DB, rule repo.DahuaEventRule, arg repo.UpdateDahuaEventRuleParams) error {
	arg.Code = strings.TrimSpace(arg.Code)
	if rule.Code == "" {
		arg.Code = rule.Code
	}

	if arg.Code == "" && rule.Code != "" {
		return errors.New("code cannot be empty")
	}

	return db.UpdateDahuaEventRule(ctx, arg)
}

func DeleteEventRule(ctx context.Context, db repo.DB, rule repo.DahuaEventRule) error {
	if rule.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.DeleteDahuaEventRule(ctx, rule.ID)
}
