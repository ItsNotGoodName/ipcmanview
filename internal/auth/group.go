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

func UpdateGroup(ctx context.Context, db repo.DB, dbModel repo.Group, arg UpdateGroupParams) (int64, error) {
	model := groupFrom(dbModel)

	model.Name = arg.Name
	model.Description = arg.Description
	model.normalize()

	if err := core.Validate.Struct(model); err != nil {
		return 0, err
	}

	return db.AuthUpdateGroup(ctx, repo.AuthUpdateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		UpdatedAt:   types.NewTime(time.Now()),
		ID:          dbModel.ID,
	})
}

func DeleteGroup(ctx context.Context, db repo.DB, id int64) error {
	return db.AuthDeleteGroup(ctx, id)
}

func UpdateGroupDisable(ctx context.Context, db repo.DB, userID int64, disable bool) error {
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
