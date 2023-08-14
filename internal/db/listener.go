package db

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type handlerFunc func(ctx context.Context, notification *pgconn.Notification, qes qes.Querier) error

type listener struct {
	config   *pgx.ConnConfig
	backlog  func(ctx context.Context, db qes.Querier) error
	handlers map[string]handlerFunc
}

func (l listener) Serve(ctx context.Context) error {
	conn, err := pgx.ConnectConfig(ctx, l.config)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	for channel := range l.handlers {
		_, err := conn.Exec(ctx, "LISTEN "+pgx.Identifier{channel}.Sanitize())
		if err != nil {
			return fmt.Errorf("LISTEN %q: %w", channel, err)
		}
	}

	if err := l.backlog(ctx, conn); err != nil {
		return err
	}

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			return fmt.Errorf("waiting for notification: %w", err)
		}

		if handler, ok := l.handlers[notification.Channel]; ok {
			err := handler(ctx, notification, conn)
			if err != nil {
				log.Err(err).Str("channel", notification.Channel).Msg("Failed to handle notification")
			}
		} else {
			log.Error().Str("channel", notification.Channel).Msg("Missing handler")
		}
	}
}
