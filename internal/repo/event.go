package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func (q *Queries) CreateEvent(ctx context.Context, action models.EventAction, data any) (int64, error) {
	actor := core.UseActor(ctx)
	b, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return q.createEvent(ctx, createEventParams{
		Action: action,
		Data:   types.NewJSON(b),
		UserID: sql.NullInt64{
			Int64: actor.UserID,
			Valid: actor.Type == core.ActorTypeUser,
		},
		Actor:     actor.Type,
		CreatedAt: types.NewTime(time.Now()),
	})
}
