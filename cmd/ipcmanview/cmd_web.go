package main

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/pubsub"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
	webserver "github.com/ItsNotGoodName/ipcmanview/internal/web/server"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/thejerf/suture/v4"
)

type WebCmd struct {
	HTTPHost string `env:"HTTP_HOST"`
	HTTPPort string `env:"HTTP_PORT" default:"8080"`
	DBPath   string `env:"DB_PATH" default:"sqlite.db"`
}

func (c *WebCmd) Run(ctx *Context) error {
	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// Bus
	dahuaBus := dahua.NewBus()

	// Database
	db, err := useDB(c.DBPath)
	if err != nil {
		return err
	}

	// Stores
	dahuaCameraStore := webdahua.NewDahuaCameraStore(db)
	dahuaStore := dahua.NewStore()
	super.Add(dahuaStore)
	eventWorkerStore :=
		dahua.NewEventWorkerStore(super,
			webdahua.NewDahuaEventHooksProxy(dahuaBus, db))
	dahua.RegisterEventBus(eventWorkerStore, dahuaBus)
	if err := dahua.Bootstrap(ctx, dahuaCameraStore, dahuaStore, eventWorkerStore); err != nil {
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
	httpServer := http.NewServer(httpRouter, c.HTTPHost+":"+c.HTTPPort)
	super.Add(httpServer)

	return super.Serve(ctx)
}
