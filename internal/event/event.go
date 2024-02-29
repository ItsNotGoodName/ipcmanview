package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/rs/zerolog/log"
)

func NewEventBuilder[T any](action string) EventBuilder[T] {
	return EventBuilder[T]{
		action: action,
	}
}

type EventBuilder[T any] struct {
	action string
	data   T
}

func (e EventBuilder[T]) Create(data T) Event {
	b, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("This should not have happend")
	}

	return Event{
		Action: e.action,
		Data:   b,
	}
}

type Event struct {
	Action string
	Data   []byte
}

func createEvent(ctx context.Context, db sqlite.DBTx, evt Event) error {
	actor := core.UseActor(ctx)

	_, err := db.CreateEvent(ctx, repo.CreateEventParams{
		Action: evt.Action,
		Data:   types.NewJSON(evt.Data),
		UserID: sql.NullInt64{
			Int64: actor.UserID,
			Valid: actor.Type == core.ActorTypeUser,
		},
		Actor:     actor.Type,
		CreatedAt: types.NewTime(time.Now()),
	})
	return err
}

func CreateEvent(ctx context.Context, db sqlite.DB, evt Event) error {
	return createEvent(ctx, db.C(), evt)
}

func CreateEventAndCommit(ctx context.Context, tx sqlite.Tx, evt Event) error {
	err := createEvent(ctx, tx.C(), evt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
