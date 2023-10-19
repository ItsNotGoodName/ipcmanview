package main

import (
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/db"
	"github.com/ItsNotGoodName/ipcmanview/migrations"
	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/server"
	"github.com/ItsNotGoodName/ipcmanview/server/api"
	"github.com/ItsNotGoodName/ipcmanview/server/rpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx, shutdown := interrupt.Context()
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

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// Bus
	bus := db.NewBusFromPool(pool)
	super.Add(bus)

	// Dahua supervisor
	dahuaSuper := dahua.NewSupervisor(pool)
	dahuaSuper.Register(bus)
	super.Add(dahuaSuper)

	// Dahua scan supervisor
	dahuaScanSuper := dahua.NewScanSupervisor(pool, dahuaSuper, 5)
	super.Add(dahuaScanSuper)

	jwt := auth.NewJWTAuthenticator(auth.JWTRandomSecret())
	authService := rpc.NewAuthService(pool, jwt)

	userService := rpc.NewUserService(pool)

	// HTTP/webrpc
	apiHandler := api.NewHandler(dahuaSuper)
	dahuaService := rpc.NewDahuaService(pool, dahuaSuper, dahuaScanSuper)
	router := server.Router(
		jwt,
		apiHandler,
		authService,
		userService,
		dahuaService,
	)
	server := server.New(router, ":8080")
	super.Add(server)

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
