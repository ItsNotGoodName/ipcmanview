package main

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/webserver"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared
	HTTPHost     string     `env:"HTTP_HOST" help:"HTTP host to listen on."`
	HTTPPort     string     `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	MQTTAddress  string     `env:"MQTT_ADDRESS" help:"MQTT broker to publish events."`
	MQTTTopic    mqtt.Topic `env:"MQTT_PREFIX" default:"ipcmanview" help:"MQTT broker topic."`
	MQTTUsername string     `env:"MQTT_USERNAME" help:"MQTT broker username."`
	MQTTPassword string     `env:"MQTT_PASSWORD" help:"MQTT broker password."`
	MQTTHa       bool       `env:"MQTT_HA" help:"Enable HomeAssistant MQTT discovery."`
	MQTTHaTopic  mqtt.Topic `env:"MQTT_HA_TOPIC" default:"homeassistant" help:"HomeAssistant MQTT discovery topic."`
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

	dahuaWorkerStore := dahuacore.NewWorkerStore(super, dahuacore.DefaultWorkerBuilder(dahuaEventHooks, bus, dahuaStore))
	dahuaWorkerStore.Register(bus)
	super.Add(sutureext.NewServiceFunc("dahua.WorkerStore", sutureext.OneShotFunc(func(ctx context.Context) error {
		return dahuaWorkerStore.Bootstrap(ctx, dahuaRepo, dahuaStore)
	})))

	bus.OnEventDahuaCoaxialStatus(func(ctx context.Context, event models.EventDahuaCoaxialStatus) error {
		fmt.Println(event)
		return nil
	})

	dahuaFileStore, err := c.useDahuaFileStore()
	if err != nil {
		return err
	}

	// MQTT
	if c.MQTTAddress != "" {
		mqttConn := mqtt.NewConn(c.MQTTTopic, c.MQTTAddress, c.MQTTUsername, c.MQTTPassword)
		super.Add(mqttConn)

		dahuaMQTTConn := dahuamqtt.NewConn(mqttConn, db, dahuaStore, c.MQTTHa, c.MQTTHaTopic)
		dahuaMQTTConn.Register(bus)
		super.Add(dahuaMQTTConn)
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
