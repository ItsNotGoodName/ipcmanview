package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/config"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuamqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuasmtp"
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
	HTTPRedirect           bool       `env:"HTTP_REDIRECT" default:"true" help:"Redirect HTTP to HTTPS."`
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
	MediamtxPathTemplate   string     `env:"MEDIAMTX_PATH_TEMPLATE" default:"ipcmanview_dahua_{{.DeviceID}}_{{.Channel}}_{{.Subtype}}" help:"Template for generating MediaMTX paths."`
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

	// Init
	auth.Init(auth.App{
		DB:                   db,
		Hub:                  hub,
		TouchSessionThrottle: auth.NewTouchSessionThrottle(),
	})
	dahua.Init(dahua.App{
		DB:         db,
		Hub:        hub,
		AFS:        dahuaAFS,
		Store:      dahuaStore,
		ScanLocker: dahua.NewScanLocker(),
	})

	// MediaMTX
	mediamtxConfig, err := mediamtx.NewConfig(c.MediamtxHost, c.MediamtxPathTemplate, c.MediamtxStreamProtocol, int(c.MediamtxWebrtcPort), int(c.MediamtxHLSPort))
	if err != nil {
		return err
	}

	// TODO: move this
	hub.OnDahuaDeviceUpdated("DEBUG", func(ctx context.Context, event bus.DahuaDeviceUpdated) error {
		client, err := mediamtx.NewClient("http://" + c.MediamtxHost + ":9997")
		if err != nil {
			return err
		}

		device, err := dahua.GetDevice(ctx, dahua.GetDeviceFilter{ID: event.DeviceID})
		if err != nil {
			return err
		}

		streams, err := db.C().DahuaListStreamsByDevice(ctx, event.DeviceID)
		if err != nil {
			return err
		}

		for _, stream := range streams {
			name := mediamtxConfig.DahuaEmbedPath(stream)
			rtspURL := dahua.GetLiveRTSPURL(dahua.GetLiveRTSPURLParams{
				Username: device.Username,
				Password: device.Password,
				Host:     device.Ip,
				Port:     554,
				Channel:  int(stream.Channel),
				Subtype:  int(stream.Subtype),
			})
			rtspTransport := "tcp"
			pathConf := mediamtx.PathConf{
				Source:        &rtspURL,
				RtspTransport: &rtspTransport,
			}

			rsp, err := client.ConfigPathsGet(ctx, name)
			if err != nil {
				return err
			}
			res, err := mediamtx.ParseConfigPathsGetResponse(rsp)
			if err != nil {
				return err
			}

			switch res.StatusCode() {
			case http.StatusOK:
				rsp, err := client.ConfigPathsPatch(ctx, name, pathConf)
				if err != nil {
					return err
				}
				rsp.Body.Close()
			case http.StatusNotFound, http.StatusInternalServerError:
				rsp, err := client.ConfigPathsAdd(ctx, name, pathConf)
				if err != nil {
					return err
				}
				rsp.Body.Close()
			default:
				return errors.New(string(res.Body))
			}
		}

		return nil
	})

	// Dahua
	if err := dahua.Normalize(ctx); err != nil {
		return err
	}

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

	dahua.RegisterStreams()

	// super.Add(dahua.NewAferoService(db, dahuaAFS))
	//
	// super.Add(dahua.NewFileService(db, dahuaAFS, dahuaStore))

	// MQTT
	if c.MQTTAddress != "" {
		mqttConn := mqtt.NewConn(c.MQTTTopic, c.MQTTAddress, c.MQTTUsername, c.MQTTPassword)
		super.Add(mqttConn)

		super.Add(dahuamqtt.NewConn(mqttConn, c.MQTTHa, c.MQTTHaTopic).Register(hub))
	}

	// SMTP
	super.Add(dahuasmtp.NewServer(dahuasmtp.NewBackend(), core.Address(c.SMTPHost, int(c.SMTPPort))))

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
		Register(rpc.NewPublicServer(rpcserver.NewPublic(configProvider), rpcLogger)).
		Register(rpc.NewUserServer(rpcserver.NewUser(configProvider, mediamtxConfig), rpcLogger, rpcserver.RequireAuthSession())).
		Register(rpc.NewAdminServer(rpcserver.NewAdmin(configProvider, db), rpcLogger, rpcserver.RequireAdminAuthSession()))

	// HTTP server
	httpServer := httpRouter
	if c.HTTPRedirect {
		httpServer = server.NewHTTPRedirect(strconv.Itoa(int(c.HTTPSPort)))
	}
	super.Add(server.NewHTTPServer(
		httpServer,
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
