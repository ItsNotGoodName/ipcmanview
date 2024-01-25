package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
)

func normalizeGroup(arg *models.Group) {
	arg.Name = strings.TrimSpace(arg.Name)
	arg.Description = strings.TrimSpace(arg.Description)
}

func CreateGroup(ctx context.Context, db repo.DB, arg models.Group) (int64, error) {
	normalizeGroup(&arg)

	if err := validate.Validate.Struct(arg); err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	return db.CreateGroup(ctx, repo.CreateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func UpdateGroup(ctx context.Context, db repo.DB, arg models.Group, newPassword string) (int64, error) {
	normalizeGroup(&arg)

	if err := validate.Validate.Struct(arg); err != nil {
		return 0, err
	}

	return db.UpdateGroup(ctx, repo.UpdateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		UpdatedAt:   types.NewTime(time.Now()),
		ID:          arg.ID,
	})
}

func DeleteGroup(ctx context.Context, db repo.DB, id int64) error {
	return db.DeleteGroup(ctx, id)
}
