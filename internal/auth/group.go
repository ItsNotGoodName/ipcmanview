package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func groupFrom(v repo.Group) group {
	return group{
		Name:        v.Name,
		Description: v.Description,
	}
}

type group struct {
	Name        string `validate:"gte=3,lte=64"`
	Description string `validate:"lte=1024"`
}

func (g *group) normalize() {
	g.Name = strings.TrimSpace(g.Name)
	g.Description = strings.TrimSpace(g.Description)
}

type CreateGroupParams struct {
	Name        string
	Description string
}

func CreateGroup(ctx context.Context, db repo.DB, arg CreateGroupParams) (int64, error) {
	if err := core.Admin(ctx); err != nil {
		return 0, err
	}

	model := group{
		Name:        arg.Name,
		Description: arg.Description,
	}
	model.normalize()

	if err := core.Validate.Struct(model); err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	return db.AuthCreateGroup(ctx, repo.AuthCreateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

type UpdateGroupParams struct {
	Name        string
	Description string
}

func UpdateGroup(ctx context.Context, db repo.DB, dbModel repo.Group, arg UpdateGroupParams) error {
	if err := core.Admin(ctx); err != nil {
		return err
	}

	model := groupFrom(dbModel)

	// Mutate
	model.Name = arg.Name
	model.Description = arg.Description
	model.normalize()

	if err := core.Validate.Struct(model); err != nil {
		return err
	}

	_, err := db.AuthUpdateGroup(ctx, repo.AuthUpdateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		UpdatedAt:   types.NewTime(time.Now()),
		ID:          dbModel.ID,
	})
	return err
}

func DeleteGroup(ctx context.Context, db repo.DB, id int64) error {
	if err := core.Admin(ctx); err != nil {
		return err
	}
	return db.AuthDeleteGroup(ctx, id)
}

func UpdateGroupDisable(ctx context.Context, db repo.DB, userID int64, disable bool) error {
	if err := core.Admin(ctx); err != nil {
		return err
	}

	if disable {
		_, err := db.AuthUpdateGroupDisabledAt(ctx, repo.AuthUpdateGroupDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         userID,
		})
		return err
	}
	_, err := db.AuthUpdateGroupDisabledAt(ctx, repo.AuthUpdateGroupDisabledAtParams{
		DisabledAt: types.NullTime{},
		ID:         userID,
	})
	return err
}
