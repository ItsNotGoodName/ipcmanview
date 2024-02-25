package main

import (
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/config"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuasmtp"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/mediamtx"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	"github.com/ItsNotGoodName/ipcmanview/internal/server"
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
	HTTPSPort              uint16     `env:"HTTPS_PORT" default:"8443" help:"HTTPS port to listen on."`
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
	if err := c.init(); err != nil {
		return err
	}

	configProvider, err := config.NewProvider(c.useConfigFilePath())
	if err != nil {
		return err
	}

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// Certificate
	cert, err := c.useCert()
	if err != nil {
		return err
	}

	// // Secret
	// if _, err := c.useSecret(); err != nil {
	// 	return err
	// }

	// Database
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	// Pub sub
	pub := pubsub.NewPub()

	// Event bus
	bus := event.NewBus().Register(pub)
	super.Add(bus)

	// Event queue
	super.Add(event.NewQueue(db, bus))

	// MediaMTX
	mediamtxConfig, err := mediamtx.NewConfig(c.MediamtxHost, c.MediamtxPathTemplate, c.MediamtxStreamProtocol, int(c.MediamtxWebrtcPort), int(c.MediamtxHLSPort))
	if err != nil {
		return err
	}

	// Dahua

	dahuaAFS, err := c.useDahuaAFS()
	if err != nil {
		return err
	}

	dahuaStore := dahua.
		NewStore(db).
		Register(bus)
	defer dahuaStore.Close()
	super.Add(dahuaStore)

	if err := dahua.
		NewWorkerManager(super, dahua.DefaultWorkerFactory(bus, pub, db, dahuaStore, dahua.NewScanLockStore(), dahua.NewDefaultEventHooks(bus, db))).
		Register(bus, db).
		Bootstrap(ctx, db, dahuaStore); err != nil {
		return err
	}

	dahua.RegisterStreams(bus, db, dahuaStore)

	super.Add(dahua.NewAferoService(db, dahuaAFS))

	super.Add(dahua.NewFileService(db, dahuaAFS, dahuaStore))

	// MQTT
	if c.MQTTAddress != "" {
		mqttConn := mqtt.NewConn(c.MQTTTopic, c.MQTTAddress, c.MQTTUsername, c.MQTTPassword)
		super.Add(mqttConn)

		super.Add(dahuamqtt.NewConn(mqttConn, db, dahuaStore, c.MQTTHa, c.MQTTHaTopic).Register(bus))
	}

	// SMTP
	super.Add(dahuasmtp.
		NewServer(dahuasmtp.
			NewBackend(dahuasmtp.NewApp(db, bus, dahuaAFS)), core.Address(c.SMTPHost, int(c.SMTPPort))))

	// HTTP router
	httpRouter := server.NewHTTPRouter(web.RouteAssets)

	// HTTP middleware
	httpRouter.Use(web.FS(api.Route, rpcserver.Route))
	httpRouter.Use(api.SessionMiddleware(db))
	httpRouter.Use(api.ActorMiddleware())

	// API
	api.
		NewServer(pub, db, bus, dahuaStore, dahuaAFS, mediamtxConfig.URL()).
		RegisterSession(httpRouter.Group(api.Route)).
		Register(httpRouter.Group(api.Route, api.RequireAuthMiddleware()))

	// RPC
	rpcLogger := rpcserver.Logger()
	rpcserver.
		NewServer(httpRouter).
		Register(rpc.NewHelloWorldServer(&rpcserver.HelloWorld{}, rpcLogger)).
		Register(rpc.NewPublicServer(rpcserver.NewPublic(configProvider, db), rpcLogger)).
		Register(rpc.NewUserServer(rpcserver.NewUser(configProvider, db, bus, dahuaStore, mediamtxConfig), rpcLogger, rpcserver.RequireAuthSession())).
		Register(rpc.NewAdminServer(rpcserver.NewAdmin(configProvider, db, bus), rpcLogger, rpcserver.RequireAdminAuthSession()))

	// HTTP server
	super.Add(server.NewHTTPServer(
		server.NewHTTPRedirect(strconv.Itoa(int(c.HTTPSPort))),
		core.Address(c.HTTPHost, int(c.HTTPPort)),
		nil,
	))

	// HTTPS server
	super.Add(server.NewHTTPServer(
		httpRouter,
		core.Address(c.HTTPHost, int(c.HTTPSPort)),
		&cert,
	))

	return super.Serve(ctx)
}
