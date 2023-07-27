package sandbox

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func notifyStart(ctx context.Context, pool *pgxpool.Pool) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Err(err).Msg("Failed get connection")
	}
	defer conn.Release()

	for _ = range []int{2, 3, 4} {
		_, err := conn.Exec(ctx, "select pg_notify('order_progress_event', 'Hello world!');")
		if err != nil {
			log.Err(err).Msg("Failed to notify")
		}
		time.Sleep(1 * time.Second)
	}
}

func Notify(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	go notifyStart(ctx, pool)

	_, err = conn.Conn().Exec(ctx, "LISTEN order_progress_event")
	if err != nil {
		return err
	}
	for {
		notify, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			return err
		}

		fmt.Println(notify.Channel, notify.Payload)
	}

	return nil
}
