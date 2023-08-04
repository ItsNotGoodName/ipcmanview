package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmango/pkg/background"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Background(bus *Bus, pool *pgxpool.Pool) background.Function {
	return background.NewFunction(background.BlockingContext, func(ctx context.Context) {
		for {
			err := Start(ctx, pool, bus)
			if errors.Is(err, context.Canceled) {
				return
			}

			log.Err(err).Msg("Event bus encountered, retrying in 15 seconds")

			if wait(ctx, 15*time.Second) {
				return
			}
		}
	})
}

func Start(ctx context.Context, pool *pgxpool.Pool, bus *Bus) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	bus.handleConnect(ctx)

	{
		batch := &pgx.Batch{}
		for _, channel := range channels {
			batch.Queue(fmt.Sprintf(`LISTEN "%s"`, channel))
		}

		br := conn.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return err
		}
	}

	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return err
		}

		bus.handle(ctx, notification)
	}
}

// wait for the duration and return false when done or return true when context is done.
func wait(ctx context.Context, duration time.Duration) bool {
	t := time.NewTicker(duration)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return true
	case <-t.C:
		return false
	}
}
