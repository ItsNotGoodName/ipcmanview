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

type CmdWeb struct {
	Shared
	HTTPHost    string `env:"HTTP_HOST" help:"HTTP host to listen on."`
	HTTPPort    string `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	MQTTDisable bool   `env:"MQTT_DISABLE" help:"Disable MQTT server."`
	MQTTHost    string `env:"MQTT_HOST" help:"MQTT host to listen on."`
	MQTTPort    string `env:"MQTT_PORT" default:"1883" help:"MQTT port to listen on."`
}

func (c *CmdWeb) Run(ctx *Context) error {
	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// Database
	db, err := useDB(c.DBPath)
	if err != nil {
		return err
	}

	// MQTT
	pub, err := pubsub.NewPub(!c.MQTTDisable, c.MQTTHost+":"+c.MQTTPort)
	if err != nil {
		return err
	}
	super.Add(pub)

	dahuaBus := dahua.NewBus()
	pub.Register(dahuaBus)

	dahuaCameraStore := webdahua.NewDahuaCameraStore(db)

	dahuaStore := dahua.NewStore()
	super.Add(dahuaStore)
	dahuaStore.Register(dahuaBus)

	eventWorkerStore := dahua.NewEventWorkerStore(super, webdahua.NewDahuaEventHooksProxy(dahuaBus, db))
	eventWorkerStore.Register(dahuaBus)

	if err := dahua.Bootstrap(ctx, dahuaCameraStore, dahuaStore, eventWorkerStore); err != nil {
		return err
	}

	// HTTP Router
	httpRouter := http.NewRouter()
	if err := webserver.RegisterRenderer(httpRouter); err != nil {
		return err
	}

	// HTTP Middleware
	webserver.RegisterMiddleware(httpRouter)

	// HTTP API
	apiDahuaServer := api.NewDahuaServer(pub, dahuaStore, dahuaCameraStore)
	api.RegisterDahuaRoutes(httpRouter, apiDahuaServer)

	// HTTP Web
	webServer := webserver.New(db, pub, dahuaStore, dahuaBus)
	webserver.RegisterRoutes(httpRouter, webServer)

	// HTTP Server
	httpServer := http.NewServer(httpRouter, c.HTTPHost+":"+c.HTTPPort)
	super.Add(httpServer)

	return super.Serve(ctx)
}
