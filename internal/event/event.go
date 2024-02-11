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

func createEvent(ctx context.Context, db sqlite.DBTx, bus *Bus, action models.EventAction, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	actor := core.UseActor(ctx)

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
	return err
}

func CreateEvent(ctx context.Context, db sqlite.DB, bus *Bus, action models.EventAction, data any) error {
	err := createEvent(ctx, db.C(), bus, action, data)
	if err != nil {
		return err
	}

	bus.EventQueued(EventQueued{})
	return nil
}

func CreateEventTx(ctx context.Context, tx sqlite.Tx, bus *Bus, action models.EventAction, data any) error {
	err := createEvent(ctx, tx.C(), bus, action, action)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	bus.EventQueued(EventQueued{})
	return nil
}
