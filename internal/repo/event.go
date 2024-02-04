package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func (q *Queries) CreateEvent(ctx context.Context, action models.EventAction, slug any) (int64, error) {
	actor := core.UseActor(ctx)
	return q.createEvent(ctx, createEventParams{
		Action: action,
		Slug:   fmt.Sprintf("%v", slug),
		UserID: sql.NullInt64{
			Int64: actor.UserID,
			Valid: actor.Type == core.ActorTypeUser,
		},
		Actor:     actor.Type,
		CreatedAt: types.NewTime(time.Now()),
	})
}
