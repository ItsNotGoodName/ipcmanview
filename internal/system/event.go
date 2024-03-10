package system

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

func CreateEvent(ctx context.Context, db sqlite.DBTx, event Event) error {
	actor := core.UseActor(ctx)

	_, err := db.CreateEvent(ctx, repo.CreateEventParams{
		Action: event.Action,
		Data:   types.NewJSON(event.Data),
		UserID: sql.NullInt64{
			Int64: actor.UserID,
			Valid: actor.Type == core.ActorTypeUser,
		},
		Actor:     actor.Type,
		CreatedAt: types.NewTime(time.Now()),
	})
	return err
}
