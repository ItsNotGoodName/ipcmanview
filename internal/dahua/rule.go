package dahua

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateEventDefaultRule(ctx context.Context, db repo.DB, arg repo.CreateDahuaEventDefaultRuleParams) error {
	if arg.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.CreateDahuaEventDefaultRule(ctx, arg)
}

func UpdateEventDefaultRule(ctx context.Context, db repo.DB, arg repo.UpdateDahuaEventDefaultRuleParams) error {
	rule, err := db.GetDahuaEventDefaultRule(ctx, arg.ID)
	if err != nil {
		return err
	}
	if rule.Code == "" {
		arg.Code = rule.Code
	}
	return db.UpdateDahuaEventDefaultRule(ctx, arg)
}

func DeleteEventDefaultRule(ctx context.Context, db repo.DB, rule repo.DahuaEventDefaultRule) error {
	if rule.Code == "" {
		return errors.New("code cannot be empty")
	}
	return db.DeleteDahuaEventDefaultRule(ctx, rule.ID)
}
