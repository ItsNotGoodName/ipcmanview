package main

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuasmtp"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/rpc"

	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared
	HTTPHost               string     `env:"HTTP_HOST" help:"HTTP host to listen on (e.g. \"127.0.0.1\")."`
	HTTPPort               uint16     `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	SMTPHost               string     `env:"SMTP_HOST" help:"SMTP host to listen on (e.g. \"127.0.0.1\")."`
	SMTPPort               uint16     `env:"SMTP_PORT" default:"1025" help:"SMTP port to listen on."`
	MQTTAddress            string     `env:"MQTT_ADDRESS" help:"MQTT server address (e.g. \"mqtt://192.168.1.20:1883\")."`
	MQTTTopic              mqtt.Topic `env:"MQTT_PREFIX" default:"ipcmanview" help:"MQTT server topic to publish messages."`
	MQTTUsername           string     `env:"MQTT_USERNAME" help:"MQTT server username for authentication."`
	MQTTPassword           string     `env:"MQTT_PASSWORD" help:"MQTT server password for authentication."`
	MQTTHa                 bool       `env:"MQTT_HA" help:"Enable Home Assistant MQTT discovery."`
	MQTTHaTopic            mqtt.Topic `env:"MQTT_HA_TOPIC" default:"homeassistant" help:"Home Assistant MQTT discover topic."`
	MediamtxHost           string     `env:"MEDIAMTX_HOST" help:"MediaMTX host address (e.g. \"192.168.1.20\")."`
	MediamtxWebrtcPort     uint16     `env:"MEDIAMTX_WEBRTC_PORT" default:"8889" help:"MediaMTX WebRTC port."`
	MediamtxHLSPort        uint16     `env:"MEDIAMTX_HLS_PORT" default:"8888" help:"MediaMTX HLS port."`
	MediamtxPathTemplate   string     `env:"MEDIAMTX_PATH_TEMPLATE" help:"Template for generating MediaMTX paths (e.g. \"ipcmanview_dahua_{{.DeviceID}}_{{.Channel}}_{{.Subtype}}\")."`
	MediamtxStreamProtocol string     `env:"MEDIAMTX_STREAM_PROTOCOL" default:"webrtc" enum:"webrtc,hls" help:"MediaMTX stream protocol."`
}

func (c *CmdServe) Run(ctx *Context) error {
	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// Secret
	_, err := c.useSecret()
	if err != nil {
		return err
	}

	// Database
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	// Pub sub
	pub := pubsub.NewPub()
	super.Add(pub)

	// Bus
	bus := event.NewBus().Register(pub)
	super.Add(bus)

	// MediaMTX
	// mediamtxConfig, err := mediamtx.NewConfig(c.MediamtxHost, c.MediamtxPathTemplate, c.MediamtxStreamProtocol, int(c.MediamtxWebrtcPort), int(c.MediamtxHLSPort))
	// if err != nil {
	// 	return err
	// }

	// Dahua

	dahuaAFS, err := c.useDahuaAFS()
	if err != nil {
		return err
	}

	dahuaStore := dahua.
		NewStore().
		Register(bus)
	defer dahuaStore.Close()
	super.Add(dahuaStore)

	dahuaScanLockStore := dahua.NewScanLockStore()

	dahuaWorkerStore := dahua.
		NewWorkerStore(super, dahua.DefaultWorkerFactory(bus, pub, db, dahuaStore, dahuaScanLockStore, dahua.NewDefaultEventHooks(bus, db))).
		Register(bus)
	if err := dahuaWorkerStore.Bootstrap(ctx, db, dahuaStore); err != nil {
		return err
	}

	dahua.RegisterStreams(bus, db, dahuaStore)

	super.Add(dahua.NewAferoService(db, dahuaAFS))

	dahuaFileService := dahua.NewFileService(db, dahuaAFS, dahuaStore)
	super.Add(dahuaFileService)

	// MQTT
	if c.MQTTAddress != "" {
		mqttConn := mqtt.NewConn(c.MQTTTopic, c.MQTTAddress, c.MQTTUsername, c.MQTTPassword)
		super.Add(mqttConn)

		dahuaMQTTConn := dahuamqtt.NewConn(mqttConn, db, dahuaStore, c.MQTTHa, c.MQTTHaTopic)
		dahuaMQTTConn.Register(bus)
		super.Add(dahuaMQTTConn)
	}

	// SMTP
	dahuaSMTPApp := dahuasmtp.NewApp(db, dahuaAFS)
	dahuaSMTPBackend := dahuasmtp.NewBackend(dahuaSMTPApp)
	dahuaSMTPServer := dahuasmtp.NewServer(dahuaSMTPBackend, core.Address(c.SMTPHost, int(c.SMTPPort)))
	super.Add(dahuaSMTPServer)

	// HTTP router
	httpRouter := http.NewRouter()

	// HTTP middleware
	httpRouter.Use(web.FS(api.Route, rpcserver.Route))
	httpRouter.Use(api.SessionMiddleware(db))

	// API
	api.
		NewServer(pub, db, dahuaStore, dahuaAFS).
		Register(httpRouter)

	// RPC
	rpcLogger := rpcserver.Logger()
	rpcserver.NewServer(httpRouter).
		Register(rpc.NewHelloWorldServer(&rpcserver.HelloWorld{}, rpcLogger)).
		Register(rpc.NewPublicServer(rpcserver.NewPublic(db), rpcLogger)).
		Register(rpc.NewUserServer(rpcserver.NewUser(db, dahuaStore), rpcLogger, rpcserver.AuthSession())).
		Register(rpc.NewAdminServer(rpcserver.NewAdmin(db, bus), rpcLogger, rpcserver.AdminAuthSession()))

	// HTTP server
	httpServer := http.NewServer(httpRouter, core.Address(c.HTTPHost, int(c.HTTPPort)))
	super.Add(httpServer)

	return super.Serve(ctx)
}
