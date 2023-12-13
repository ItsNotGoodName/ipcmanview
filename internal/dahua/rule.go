package dahua

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateEventRule(ctx context.Context, db repo.DB, arg repo.CreateDahuaEventRuleParams) error {
	if arg.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.CreateDahuaEventRule(ctx, arg)
}

func UpdateEventRule(ctx context.Context, db repo.DB, arg repo.UpdateDahuaEventRuleParams) error {
	rule, err := db.GetDahuaEventRule(ctx, arg.ID)
	if err != nil {
		return err
	}
	if rule.Code == "" {
		arg.Code = rule.Code
	}
	return db.UpdateDahuaEventRule(ctx, arg)
}

func DeleteEventRule(ctx context.Context, db repo.DB, rule repo.DahuaEventRule) error {
	if rule.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.DeleteDahuaEventRule(ctx, rule.ID)
}
