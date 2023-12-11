package main

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuaweb"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/webserver"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared
	HTTPHost     string `env:"HTTP_HOST" help:"HTTP host to listen on."`
	HTTPPort     string `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	MQTTAddress  string `env:"MQTT_ADDRESS" help:"MQTT broker to publish events."`
	MQTTPrefix   string `env:"MQTT_PREFIX" default:"ipcmanview/" help:"MQTT broker topic prefix"`
	MQTTUsername string `env:"MQTT_USERNAME" help:"MQTT broker username."`
	MQTTPassword string `env:"MQTT_PASSWORD" help:"MQTT broker password."`
}

func (c *CmdServe) Run(ctx *Context) error {
	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// Database
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	pub := pubsub.NewPub()
	super.Add(pub)

	dahuaBus := dahua.NewBus()
	dahuaBus.Register(pub)

	dahuaRepo := dahuaweb.NewRepo(db)

	dahuaStore := dahua.NewStore()
	super.Add(dahuaStore)
	dahuaStore.Register(dahuaBus)

	eventWorkerStore := dahua.NewEventWorkerStore(super, dahuaweb.NewEventHooksProxy(dahuaBus, db))
	eventWorkerStore.Register(dahuaBus)

	if c.MQTTAddress != "" {
		mqttPublisher := mqtt.NewPublisher(c.MQTTPrefix, c.MQTTAddress, c.MQTTUsername, c.MQTTPassword)
		super.Add(mqttPublisher)
		mqttPublisher.Register(dahuaBus)
	}

	if err := dahua.Bootstrap(ctx, dahuaRepo, dahuaStore, eventWorkerStore); err != nil {
		return err
	}

	dahuaFileStore, err := c.useDahuaFileStore()
	if err != nil {
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
	apiServer := api.NewServer(pub, dahuaStore, dahuaRepo, dahuaFileStore)
	apiServer.RegisterDahuaRoutes(httpRouter)

	// HTTP Web
	webServer := webserver.New(db, pub, dahuaStore, dahuaBus)
	webserver.RegisterRoutes(httpRouter, webServer)

	// HTTP Server
	httpServer := http.NewServer(httpRouter, c.HTTPHost+":"+c.HTTPPort)
	super.Add(httpServer)

	return super.Serve(ctx)
}
