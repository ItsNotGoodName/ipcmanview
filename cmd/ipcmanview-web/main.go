package main

import (
	"context"
	"flag"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/config"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
	webserver "github.com/ItsNotGoodName/ipcmanview/internal/web/server"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/Rican7/lieut"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	ctx := context.Background()

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg := config.NewWeb(flags)

	app := lieut.NewSingleCommandApp(
		lieut.AppInfo{
			Name:    "ipcmanview-web",
			Version: build.Current.Version,
			Summary: "Basic web application for accessing IP Cameras.",
		},
		run(flags, cfg),
		flags,
		os.Stdout,
		os.Stderr,
	)

	code := app.Run(ctx, os.Args[1:])

	os.Exit(code)
}

func run(flags *flag.FlagSet, cfg *config.Web) lieut.Executor {
	return func(ctx context.Context, arguments []string) error {
		err := flags.Parse(arguments)
		if err != nil {
			return err
		}

		// Supervisor
		super := suture.New("root", suture.Spec{
			EventHook: sutureext.EventHook(),
		})

		// Bus
		dahuaBus := dahua.NewBus()

		// Database
		sqlDB, err := sqlite.New(cfg.DBPath)
		if err != nil {
			return err
		}
		if err := migrations.Migrate(sqlDB); err != nil {
			return err
		}
		db := sqlc.NewDB(sqlite.NewDebugDB(sqlDB))

		// Stores
		dahuaCameraStore := webdahua.NewDahuaCameraStore(db)
		dahuaStore := dahua.NewStore()
		super.Add(dahuaStore)
		eventWorkerStore :=
			dahua.NewEventWorkerStore(super,
				webdahua.NewDahuaEventHooksProxy(dahuaBus, db))
		dahua.RegisterEventBus(eventWorkerStore, dahuaBus)
		if err := core.DahuaBootstrap(ctx, dahuaCameraStore, dahuaStore, eventWorkerStore); err != nil {
			return err
		}

		pubSub := pubsub.NewPub(dahuaBus)

		// HTTP Router
		httpRouter := http.NewRouter()
		if err := webserver.RegisterRenderer(httpRouter); err != nil {
			return err
		}

		// HTTP Middleware
		webserver.RegisterMiddleware(httpRouter)

		// HTTP API
		apiDahuaServer := api.NewDahuaServer(pubSub, dahuaStore, dahuaCameraStore)
		api.RegisterDahuaRoutes(httpRouter, apiDahuaServer)

		// HTTP Web
		webServer := webserver.New(db, pubSub, dahuaStore, dahuaBus)
		webserver.RegisterRoutes(httpRouter, webServer)

		// HTTP Server
		httpServer := http.NewServer(httpRouter, cfg.HTTPHost+":"+cfg.HTTPPort)
		super.Add(httpServer)

		return super.Serve(ctx)
	}
}
