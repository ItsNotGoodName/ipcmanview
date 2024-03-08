package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func groupFrom(v repo.Group) _Group {
	return _Group{
		Name:        v.Name,
		Description: v.Description,
	}
}

type _Group struct {
	Name        string `validate:"gte=3,lte=64"`
	Description string `validate:"lte=1024"`
}

func (g *_Group) normalize() {
	g.Name = strings.TrimSpace(g.Name)
	g.Description = strings.TrimSpace(g.Description)
}

type CreateGroupParams struct {
	Name        string
	Description string
}

func CreateGroup(ctx context.Context, arg CreateGroupParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	model := _Group{
		Name:        arg.Name,
		Description: arg.Description,
	}
	model.normalize()

	if err := core.ValidateStruct(ctx, model); err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	return app.DB.C().AuthCreateGroup(ctx, repo.AuthCreateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

type UpdateGroupParams struct {
	ID          int64
	Name        string
	Description string
}

func UpdateGroup(ctx context.Context, arg UpdateGroupParams) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	dbModel, err := app.DB.C().AuthGetGroup(ctx, arg.ID)
	if err != nil {
		return err
	}
	model := groupFrom(dbModel)

	// Mutate
	model.Name = arg.Name
	model.Description = arg.Description
	model.normalize()

	if err := core.ValidateStruct(ctx, model); err != nil {
		return err
	}

	_, err = app.DB.C().AuthUpdateGroup(ctx, repo.AuthUpdateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		UpdatedAt:   types.NewTime(time.Now()),
		ID:          dbModel.ID,
	})
	return err
}

func DeleteGroup(ctx context.Context, id int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}
	return app.DB.C().AuthDeleteGroup(ctx, id)
}

func UpdateGroupDisable(ctx context.Context, userID int64, disable bool) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	_, err := app.DB.C().AuthUpdateGroupDisabledAt(ctx, repo.AuthUpdateGroupDisabledAtParams{
		DisabledAt: types.NullTime{
			Time:  types.NewTime(time.Now()),
			Valid: disable,
		},
		ID: userID,
	})
	return err
}
