package main

import (
	"context"
	"flag"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/migrations"
	"github.com/ItsNotGoodName/ipcmanview/internal/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	webcore "github.com/ItsNotGoodName/ipcmanview/internal/web/core"
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

	app := lieut.NewSingleCommandApp(
		lieut.AppInfo{
			Name:    "ipcmanview-web",
			Version: build.Version,
			Summary: "Basic web interface for accessing IP Cameras.",
		},
		run(),
		flags,
		os.Stdout,
		os.Stderr,
	)

	code := app.Run(ctx, os.Args[1:])

	os.Exit(code)
}

func run() lieut.Executor {
	return func(ctx context.Context, arguments []string) error {
		// Supervisor
		super := suture.New("root", suture.Spec{
			EventHook: sutureext.EventHook(),
		})

		// Database
		sqlDB, err := sqlite.New(os.Getenv("DB_PATH"))
		if err != nil {
			return err
		}
		sqliteDB := sqlite.NewDebugDB(sqlDB)
		if err := migrations.Migrate(sqliteDB); err != nil {
			return err
		}
		sqlcDB := sqlc.New(sqliteDB)

		// Bus
		dahuaBus := dahua.NewBus()

		// Stores
		dahuaStore := dahua.NewStore(dahuaBus)
		super.Add(dahuaStore)
		eventWorkerStore := dahua.NewEventWorkerStore(super, dahuaBus)
		dahua.RegisterEventBus(eventWorkerStore, dahuaBus)
		if err := webcore.SyncDahuaStore(ctx, sqlcDB, dahuaStore); err != nil {
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
		apiDahuaServer := api.
			NewDahuaServer(webcore.
				NewDahuaStoreProxy(dahuaStore, sqlcDB), pubSub)
		api.RegisterDahuaRoutes(httpRouter, apiDahuaServer)

		// HTTP Web
		webServer := webserver.New(sqlcDB, dahuaStore, pubSub)
		webserver.RegisterRoutes(httpRouter, webServer)

		// HTTP Server
		httpServer := http.NewServer(httpRouter, ":8080")
		super.Add(httpServer)

		return super.Serve(ctx)
	}
}
