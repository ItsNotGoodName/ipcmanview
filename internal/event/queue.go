package event

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func NewQueue(db repo.DB, bus *Bus) Queue {
	q := Queue{
		db:    db,
		bus:   bus,
		check: make(chan struct{}, 1),
	}

	q.bus.OnEventQueued(func(ctx context.Context, evt EventQueued) error {
		select {
		case q.check <- struct{}{}:
		default:
		}
		return nil
	})

	return q
}

type Queue struct {
	db    repo.DB
	bus   *Bus
	check chan struct{}
}

func (q Queue) Serve(ctx context.Context) error {
	cursor, err := q.db.GetEventCursor(ctx)
	if err != nil && !repo.IsNotFound(err) {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-q.check:
			for {
				event, err := q.db.NextEventByCursor(ctx, cursor)
				if err != nil {
					if repo.IsNotFound(err) {
						break
					}
					return err
				}
				cursor = event.ID

				q.bus.Event(Event{
					Event: event,
				})
			}
		}
	}
}
