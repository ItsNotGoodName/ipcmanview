package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/mediamtx"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	webserver "github.com/ItsNotGoodName/ipcmanview/internal/web/server"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuaevents"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/encode"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/rpc"

	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared
	HTTPHost             string     `env:"HTTP_HOST" help:"HTTP host to listen on (e.g. \"127.0.0.1\")."`
	HTTPPort             string     `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	MQTTAddress          string     `env:"MQTT_ADDRESS" help:"MQTT server address (e.g. \"mqtt://example.com:1883\")."`
	MQTTTopic            mqtt.Topic `env:"MQTT_PREFIX" default:"ipcmanview" help:"MQTT server topic to publish messages."`
	MQTTUsername         string     `env:"MQTT_USERNAME" help:"MQTT server username for authentication."`
	MQTTPassword         string     `env:"MQTT_PASSWORD" help:"MQTT server password for authentication."`
	MQTTHa               bool       `env:"MQTT_HA" help:"Enable Home Assistant MQTT discovery."`
	MQTTHaTopic          mqtt.Topic `env:"MQTT_HA_TOPIC" default:"homeassistant" help:"Home Assistant MQTT discover topic."`
	MediamtxWebAddress   string     `env:"MEDIAMTX_WEB_ADDRESS" help:"MediaMTX web server address for streaming (e.g. \"http://example.com:8889\" or \"http://example.com:8888\")."`
	MediamtxPathTemplate string     `env:"MEDIAMTX_PATH_TEMPLATE" help:"Template for generating MediaMTX paths (e.g. \"ipcmanview_dahua_{{.DeviceID}}_{{.Channel}}_{{.Subtype}}\")."`
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

	// MediaMTX
	mediamtxConfig, err := mediamtx.NewConfig(c.MediamtxWebAddress, c.MediamtxPathTemplate)
	if err != nil {
		return err
	}

	// Dahua

	dahuaStore := dahua.NewStore().Register(bus)
	super.Add(dahuaStore)

	dahuaEventHooks := dahua.NewDefaultEventHooks(bus, db)
	super.Add(dahuaEventHooks)

	dahuaWorkerStore := dahua.NewWorkerStore(super, dahua.DefaultWorkerFactory(bus, pub, db, dahuaStore, dahuaEventHooks)).Register(bus)
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
	webserver.New(db, pub, bus, dahuaStore, mediamtxConfig).Register(httpRouter)

	// API
	api.NewServer(pub, db, dahuaStore, dahuaFileStore).Register(httpRouter)

	// RPC
	rpcserver.Register(httpRouter, rpc.NewHelloWorldServer(&rpcserver.HelloWorld{}))

	// HTTP server
	httpServer := http.NewServer(httpRouter, c.HTTPHost+":"+c.HTTPPort)
	super.Add(httpServer)

	// TODO: move this
	bus.OnEventDahuaEvent(func(ctx context.Context, event models.EventDahuaEvent) error {
		switch event.Event.Code {
		case dahuaevents.CodeNewFile:
			go func() error {

				dbDevice, err := db.GetDahuaDevice(ctx, event.Event.DeviceID)
				if err != nil {
					return err
				}
				client := dahuaStore.Client(ctx, dbDevice.Convert().DahuaConn)

				err = dahua.ScanLockCreateTry(ctx, db, client.Conn.ID)
				if err != nil {
					return err
				}
				cancel := dahua.ScanLockHeartbeat(ctx, db, client.Conn.ID)
				defer cancel()

				return dahua.Scan(ctx, db, client.RPC, client.Conn, models.DahuaScanTypeQuick)
			}()

			return nil
		default:
			return nil
		}
	})

	// TODO: move this
	dahuaStreamsHandle := func(ctx context.Context, device models.DahuaDeviceConn) error {
		if !device.DahuaConn.Feature.EQ(models.DahuaFeatureCamera) {
			return nil
		}

		client := dahuaStore.Client(ctx, device.DahuaConn)
		caps, err := encode.GetCaps(ctx, client.RPC, 1)
		if err != nil {
			return err
		}

		subtypes := 1
		if caps.MaxExtraStream > 0 && caps.MaxExtraStream < 10 {
			subtypes += caps.MaxExtraStream
		}

		for channelIndex, device := range caps.VideoEncodeDevices {
			names := make([]string, subtypes)
			for i, v := range device.SupportDynamicBitrate {
				if i < len(names) {
					names[i] = v.Stream
				}
			}

			args := repo.TryCreateDahuaStreamParams{
				DeviceID: client.Conn.ID,
				Channel:  int64(channelIndex + 1),
			}
			for i := 0; i < subtypes; i++ {
				args.Subtype = int64(i)
				args.Name = names[i]
				err := db.TryCreateDahuaStream(ctx, args)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
	bus.OnEventDahuaDeviceCreated(func(ctx context.Context, event models.EventDahuaDeviceCreated) error {
		return dahuaStreamsHandle(ctx, event.Device)
	})
	bus.OnEventDahuaDeviceUpdated(func(ctx context.Context, event models.EventDahuaDeviceUpdated) error {
		return dahuaStreamsHandle(ctx, event.Device)
	})

	return super.Serve(ctx)
}
