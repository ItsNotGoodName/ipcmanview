package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	webserver "github.com/ItsNotGoodName/ipcmanview/internal/web/server"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuaevents"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/rpc"

	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared
	HTTPHost     string     `env:"HTTP_HOST" help:"HTTP host to listen on."`
	HTTPPort     string     `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	MQTTAddress  string     `env:"MQTT_ADDRESS" help:"MQTT broker to publish events."`
	MQTTTopic    mqtt.Topic `env:"MQTT_PREFIX" default:"ipcmanview" help:"MQTT broker publish topic."`
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
	bus := core.NewBus().Register(pub)
	super.Add(bus)

	// Dahua

	dahuaStore := dahua.NewStore().Register(bus)
	super.Add(dahuaStore)

	dahuaEventHooks := dahua.NewDefaultEventHooks(bus, db)
	super.Add(dahuaEventHooks)

	dahuaWorkerStore := dahua.NewWorkerStore(super, dahua.DefaultWorkerBuilder(dahuaEventHooks, bus, dahuaStore, db)).Register(bus)
	if err := dahuaWorkerStore.Bootstrap(ctx, db, dahuaStore); err != nil {
		return err
	}

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

	// HTTP router
	httpRouter := http.NewRouter()
	if err := webserver.RegisterRenderer(httpRouter); err != nil {
		return err
	}

	// WEB
	webserver.New(db, pub, bus, dahuaStore).Register(httpRouter)

	// API
	api.NewServer(pub, db, dahuaStore, dahuaFileStore).Register(httpRouter)

	// RPC
	rpcserver.Register(httpRouter, rpc.NewHelloWorldServer(&rpcserver.HelloWorld{}))

	// HTTP server
	httpServer := http.NewServer(httpRouter, c.HTTPHost+":"+c.HTTPPort)
	super.Add(httpServer)

	// TODO: remove this
	bus.OnEventDahuaEvent(func(ctx context.Context, event models.EventDahuaEvent) error {
		switch event.Event.Code {
		case dahuaevents.CodeNewFile:
			go func() error {

				dbDevice, err := db.GetDahuaDevice(ctx, event.Event.DeviceID)
				if err != nil {
					return err
				}
				client := dahuaStore.Client(ctx, dbDevice.Convert().DahuaConn)

				err = dahua.ScanLockCreate(ctx, db, client.Conn.ID)
				if err != nil {
					return err
				}
				cancel := dahua.ScanLockHeartbeat(ctx, db, client.Conn.ID)
				defer cancel()

				return dahua.Scan(ctx, db, client.RPC, client.Conn, dahua.ScanTypeQuick)
			}()

			return nil
		default:
			return nil
		}
	})

	return super.Serve(ctx)
}
