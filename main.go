package main

import (
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/build"
	"github.com/ItsNotGoodName/ipcmango/migrations"
	"github.com/ItsNotGoodName/ipcmango/pkg/interrupt"
	pgxzerolog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := interrupt.Context()

	// Database
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse pgx config")
	}
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzerolog.NewLogger(log.Logger),
		LogLevel: tracelog.LogLevelDebug,
	}
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create connection pool")
	}
	defer pool.Close()
	if err := migrations.Migrate(ctx, pool); err != nil {
		log.Err(err).Msg("Failed to migrate database")
	}
}

var (
	builtBy    = "unknown"
	commit     = ""
	date       = ""
	version    = "dev"
	repoURL    = "https://github.com/ItsNotGoodName/smtpbridge"
	releaseURL = ""
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	build.Current = build.Build{
		BuiltBy:    builtBy,
		Commit:     commit,
		Date:       date,
		Version:    version,
		RepoURL:    repoURL,
		ReleaseURL: releaseURL,
	}
}
