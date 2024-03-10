package main

import (
	"context"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuasmtp"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuatasks"
	"github.com/ItsNotGoodName/ipcmanview/internal/mediamtx"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	"github.com/ItsNotGoodName/ipcmanview/internal/server"
	"github.com/ItsNotGoodName/ipcmanview/internal/squeuel"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/thejerf/suture/v4"
)

type CmdServe struct {
	Shared

	HttpHost     string `env:"HTTP_HOST" help:"HTTP(S) host to listen on (e.g. \"127.0.0.1\")."`
	HttpPort     uint16 `env:"HTTP_PORT" default:"8080" help:"HTTP port to listen on."`
	HttpsPort    uint16 `env:"HTTPS_PORT" default:"8443" help:"HTTPS port to listen on."`
	HttpRedirect bool   `env:"HTTP_REDIRECT" default:"true" negatable:"" help:"Redirect HTTP to HTTPS."`

	SmtpHost string `env:"SMTP_HOST" help:"SMTP host to listen on (e.g. \"127.0.0.1\")."`
	SmtpPort uint16 `env:"SMTP_PORT" default:"1025" help:"SMTP port to listen on."`

	MqttAddress   string     `env:"MQTT_ADDRESS" help:"MQTT server address (e.g. \"mqtt://192.168.1.20:1883\")."`
	MqttTopic     mqtt.Topic `env:"MQTT_PREFIX" default:"ipcmanview" help:"MQTT server topic to publish messages."`
	MqttUsername  string     `env:"MQTT_USERNAME" help:"MQTT server username for authentication."`
	MqttPassword  string     `env:"MQTT_PASSWORD" help:"MQTT server password for authentication."`
	MqttHass      bool       `env:"MQTT_HASS" help:"Enable Home Assistant MQTT discovery."`
	MqttHassTopic mqtt.Topic `env:"MQTT_HASS_TOPIC" default:"homeassistant" help:"Home Assistant MQTT discover topic."`

	MediamtxHost           string `env:"MEDIAMTX_HOST" help:"MediaMTX host address (e.g. \"192.168.1.20\")."`
	MediamtxApiHost        string `env:"MEDIAMTX_API_HOST" help:"MediaMTX API host (e.g. \"192.168.1.20\")."`
	MediamtxApiPort        uint16 `env:"MEDIAMTX_API_PORT" default:"9997" help:"MediaMTX API port."`
	MediamtxWebrtcPort     uint16 `env:"MEDIAMTX_WEBRTC_PORT" default:"8889" help:"MediaMTX WebRTC port."`
	MediamtxHlsPort        uint16 `env:"MEDIAMTX_HLS_PORT" default:"8888" help:"MediaMTX HLS port."`
	MediamtxStreamProtocol string `env:"MEDIAMTX_STREAM_PROTOCOL" default:"webrtc" enum:"webrtc,hls" help:"MediaMTX stream protocol (webrtc,hls)."`
}

func (c *CmdServe) Run(ctx *Context) error {
	if err := c.init(); err != nil {
		return err
	}

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	configProvider, err := system.NewConfigProvider(c.useConfigFilePath())
	if err != nil {
		return err
	}

	// Certificate
	cert, err := c.useCert()
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

	// Bus hub
	hub := bus.NewHub(ctx).Register(pub)

	// Dahua AFS
	dahuaAFS, err := c.useDahuaAFS()
	if err != nil {
		return err
	}

	// Dahua store
	dahuaStore := dahua.NewStore().Register(hub)
	defer dahuaStore.Close()

	// MediaMTX
	mediamtxConfig, err := mediamtx.NewConfig(c.MediamtxHost, c.MediamtxStreamProtocol, int(c.MediamtxWebrtcPort), int(c.MediamtxHlsPort))
	if err != nil {
		return err
	}

	mediamtxClient, err := mediamtx.NewClient("http://" + core.Address(core.First(c.MediamtxApiHost, c.MediamtxHost), int(c.MediamtxApiPort)))
	if err != nil {
		return err
	}

	// Init
	system.Init(system.App{
		CP: configProvider,
	})
	auth.Init(auth.App{
		DB:                   db,
		Hub:                  hub,
		TouchSessionThrottle: auth.NewTouchSessionThrottle(),
	})
	dahua.Init(dahua.App{
		DB:             db,
		Hub:            hub,
		AFS:            dahuaAFS,
		Store:          dahuaStore,
		ScanLocker:     dahua.NewScanLocker(),
		MediamtxClient: mediamtxClient,
		MediamtxConfig: mediamtxConfig,
	})
	dahuatasks.Init(dahuatasks.App{
		DB:  db,
		Hub: hub,
	})

	// Dahua
	if err := dahua.Normalize(ctx); err != nil {
		return err
	}

	// Sync stream queue
	super.Add(squeuel.NewWorker(db, dahuatasks.SyncStreamTask.Queue, dahuatasks.HandleSyncStreamTask).Register(hub))
	// Push stream queue
	super.Add(squeuel.NewWorker(db, dahuatasks.PushStreamTask.Queue, dahuatasks.HandlePushStreamTask).Register(hub))

	dahuatasks.RegisterStreams()

	dahuaWorkerHooks := dahua.NewDefaultWorkerHooks()

	if err := dahua.
		NewWorkerManager(super, func(ctx context.Context, super *suture.Supervisor, conn dahua.Conn) []suture.ServiceToken {
			return []suture.ServiceToken{
				super.Add(dahua.NewQuickScanWorker(dahuaWorkerHooks, pub, conn.ID)),
				super.Add(dahua.NewCoaxialWorker(dahuaWorkerHooks, conn.ID)),
				super.Add(dahua.NewEventWorker(dahuaWorkerHooks, conn)),
			}
		}).
		Register().
		Bootstrap(ctx); err != nil {
		return err
	}

	// super.Add(dahua.NewAferoService(db, dahuaAFS))
	//
	// super.Add(dahua.NewFileService(db, dahuaAFS, dahuaStore))

	// MQTT
	if c.MqttAddress != "" {
		mqttConn := mqtt.NewConn(c.MqttTopic, c.MqttAddress, c.MqttUsername, c.MqttPassword)
		super.Add(mqttConn)

		super.Add(dahuamqtt.NewConn(mqttConn, c.MqttHass, c.MqttHassTopic).Register(hub))
	}

	// SMTP
	super.Add(dahuasmtp.NewServer(dahuasmtp.NewBackend(), core.Address(c.SmtpHost, int(c.SmtpPort))))

	// HTTP router
	httpRouter := server.NewHTTPRouter(web.RouteAssets)

	// HTTP middleware
	httpRouter.Use(web.FS(api.Route, rpcserver.Route))
	httpRouter.Use(api.SessionMiddleware())
	httpRouter.Use(api.ActorMiddleware())

	// API
	api.
		NewServer(pub, db, dahuaAFS, mediamtxConfig.URL()).
		RegisterSession(httpRouter.Group(api.Route)).
		Register(httpRouter.Group(api.Route, api.RequireAuthMiddleware()))

	// RPC
	rpcLogger := rpcserver.Logger()
	rpcserver.
		NewServer(httpRouter).
		Register(rpc.NewHelloWorldServer(&rpcserver.HelloWorld{}, rpcLogger)).
		Register(rpc.NewPublicServer(rpcserver.NewPublic(), rpcLogger)).
		Register(rpc.NewUserServer(rpcserver.NewUser(mediamtxConfig), rpcLogger, rpcserver.RequireAuthSession())).
		Register(rpc.NewAdminServer(rpcserver.NewAdmin(db), rpcLogger, rpcserver.RequireAdminAuthSession()))

	// HTTP server
	httpServer := httpRouter
	if c.HttpRedirect {
		httpServer = server.NewHTTPRedirect(strconv.Itoa(int(c.HttpsPort)))
	}
	super.Add(server.NewHTTPServer(
		httpServer,
		core.Address(c.HttpHost, int(c.HttpPort)),
		nil,
	))

	// HTTPS server
	super.Add(server.NewHTTPServer(
		httpRouter,
		core.Address(c.HttpHost, int(c.HttpsPort)),
		&cert,
	))

	return super.Serve(ctx)
}
