package squeuel

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
)

func NewWorker(db sqlite.DB, queue string, fn HandleFunc) Worker {
	return Worker{
		db:    db,
		queue: queue,
		fn:    fn,
		flagC: make(chan struct{}, 1),
	}
}

type Worker struct {
	db    sqlite.DB
	queue string
	fn    HandleFunc
	flagC chan struct{}
}

func (w Worker) String() string {
	return fmt.Sprintf("squeuel.Worker(queue=%s)", w.queue)
}

func (w Worker) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w Worker) serve(ctx context.Context) error {
	core.FlagChannel(w.flagC)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.flagC:
			for {
				more, err := Do(ctx, w.db, w.queue, w.fn)
				if err != nil {
					return err
				}
				if !more {
					break
				}
			}
		}
	}
}

func (w Worker) Flag() {
	core.FlagChannel(w.flagC)
}

func (w Worker) Register(hub *bus.Hub) Worker {
	hub.OnSqueuelEnqueued(w.String(), func(ctx context.Context, event bus.SqueuelEnqueued) error {
		if w.queue == event.Queue {
			w.Flag()
		}
		return nil
	})
	return w
}
