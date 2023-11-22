package main

import (
	"context"
	"flag"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/pubsub"
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
			Name:    "ipcmanview-gateway",
			Version: build.Version,
			Summary: "API gateway for accessing IP Cameras.",
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

		// Bus
		dahuaBus := dahua.NewBus()

		// Stores
		dahuaCameraStore := dahua.NewCameraStore()
		dahuaStore := dahua.NewStore(dahuaCameraStore)
		super.Add(dahuaStore)
		eventWorkerStore := dahua.NewEventWorkerStore(super, dahuaBus)
		dahua.RegisterEventBus(eventWorkerStore, dahuaBus)

		pubSub := pubsub.NewPub(dahuaBus)

		// HTTP Router
		httpRouter := http.NewRouter()

		// HTTP API
		apiDahuaServer := api.NewDahuaServer(pubSub, dahuaStore, dahuaCameraStore)
		api.RegisterDahuaRoutes(httpRouter, apiDahuaServer)

		// HTTP Server
		httpServer := http.NewServer(httpRouter, ":8080")
		super.Add(httpServer)

		return super.Serve(ctx)
	}
}
