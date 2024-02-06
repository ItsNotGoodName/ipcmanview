package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func CreateEvent(ctx context.Context, db sqlite.DBTx, bus *Bus, action models.EventAction, data any) error {
	actor := core.UseActor(ctx)
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = db.CreateEvent(ctx, repo.CreateEventParams{
		Action: action,
		Data:   types.NewJSON(b),
		UserID: sql.NullInt64{
			Int64: actor.UserID,
			Valid: actor.Type == core.ActorTypeUser,
		},
		Actor:     actor.Type,
		CreatedAt: types.NewTime(time.Now()),
	})
	if err != nil {
		return err
	}
	bus.EventQueued(EventQueued{})
	return err
}
