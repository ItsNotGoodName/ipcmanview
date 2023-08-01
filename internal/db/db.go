package db

import (
	"context"
	"time"

	pgxzerolog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzerolog.NewLogger(log.Logger),
		LogLevel: tracelog.LogLevelDebug,
	}
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	return pgxpool.NewWithConfig(ctx, config)
}

func NewConn(ctx context.Context, pool *pgxpool.Pool) (*pgx.Conn, error) {
	return pgx.ConnectConfig(ctx, pool.Config().ConnConfig)
}

type Context struct {
	context.Context
	Conn *pgx.Conn
}
