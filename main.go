package main

import (
	"context"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/db"
	"github.com/ItsNotGoodName/ipcmanview/migrations"
	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/sandbox"
	"github.com/ItsNotGoodName/ipcmanview/server"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx, shutdown := context.WithCancel(interrupt.Context())
	defer shutdown()

	// Database
	pool, err := db.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database connection pool")
	}
	defer pool.Close()

	// Database migrate
	if err := migrations.Migrate(ctx, pool); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// sandbox.Sandbox(ctx, pool)
	//
	// return

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// HTTP/webrpc
	http := server.NewHTTP(chi.NewRouter(), ":8080", shutdown)
	super.Add(http)

	// DEBUG
	super.Add(sandbox.NewSandbox(pool))

	super.Serve(ctx)
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
