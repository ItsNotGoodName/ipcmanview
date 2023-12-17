package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/hass"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/webserver"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared
	HTTPHost      string     `env:"HTTP_HOST" help:"HTTP host to listen on."`
	HTTPPort      string     `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	MQTTAddress   string     `env:"MQTT_ADDRESS" help:"MQTT broker to publish events."`
	MQTTTopic     mqtt.Topic `env:"MQTT_PREFIX" default:"ipcmanview" help:"MQTT broker topic."`
	MQTTUsername  string     `env:"MQTT_USERNAME" help:"MQTT broker username."`
	MQTTPassword  string     `env:"MQTT_PASSWORD" help:"MQTT broker password."`
	MQTTHass      bool       `env:"MQTT_HASS" help:"Enable HomeAssistant MQTT discovery"`
	MQTTHassTopic mqtt.Topic `env:"MQTT_HASS_TOPIC" default:"homeassistant" help:"HomeAssistant MQTT discovery topic"`
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

	// Pub sub
	pub := pubsub.NewPub()
	super.Add(pub)

	// Bus
	bus := core.NewBus()
	bus.Register(pub)
	super.Add(bus)

	// Dahua

	dahuaRepo := dahua.NewRepo(db)

	dahuaStore := dahuacore.NewStore()
	dahuaStore.Register(bus)
	super.Add(dahuaStore)

	dahuaEventHooks := dahua.NewEventHooks(bus, db)
	super.Add(dahuaEventHooks)

	dahuaEventWorkerStore := dahuacore.NewEventWorkerStore(super, dahuaEventHooks)
	dahuaEventWorkerStore.Register(bus)

	dahuaFileStore, err := c.useDahuaFileStore()
	if err != nil {
		return err
	}

	super.Add(sutureext.NewServiceFunc("dahua.bootstrap", sutureext.OneShotFunc(func(ctx context.Context) error {
		return dahuacore.Bootstrap(ctx, dahuaRepo, dahuaStore, dahuaEventWorkerStore)
	})))

	// MQTT
	if c.MQTTAddress != "" {
		mqttConn := mqtt.NewConn(c.MQTTTopic, c.MQTTAddress, c.MQTTUsername, c.MQTTPassword)
		mqtt.Register(mqttConn, bus)
		super.Add(mqttConn)

		hassConn := hass.NewConn(mqttConn, db, c.MQTTHassTopic)
		super.Add(hassConn)
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
	webServer := webserver.New(db, pub, bus, dahuaStore)
	webserver.RegisterRoutes(httpRouter, webServer)

	// HTTP Server
	httpServer := http.NewServer(httpRouter, c.HTTPHost+":"+c.HTTPPort)
	super.Add(httpServer)

	return super.Serve(ctx)
}
