package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
)

func NewGroup(v repo.Group) Group {
	return Group{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}
}

type Group struct {
	ID          int64
	Name        string `validate:"gte=3,lte=64"`
	Description string `validate:"lte=1024"`
}

func (g *Group) normalize() {
	g.Name = strings.TrimSpace(g.Name)
	g.Description = strings.TrimSpace(g.Description)
}

func CreateGroup(ctx context.Context, db repo.DB, arg Group) (int64, error) {
	arg.normalize()

	if err := validate.Validate.Struct(arg); err != nil {
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

func UpdateGroup(ctx context.Context, db repo.DB, arg Group) (int64, error) {
	arg.normalize()

	if err := validate.Validate.Struct(arg); err != nil {
		return 0, err
	}

	return db.AuthUpdateGroup(ctx, repo.AuthUpdateGroupParams{
		Name:        arg.Name,
		Description: arg.Description,
		UpdatedAt:   types.NewTime(time.Now()),
		ID:          arg.ID,
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
