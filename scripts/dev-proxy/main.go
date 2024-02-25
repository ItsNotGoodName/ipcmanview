// Package main runs the HTTP development proxy.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	"github.com/ItsNotGoodName/ipcmanview/internal/server"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/pkg/echoext"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address     string
	Certificate *server.Certificate
	Servers     []ConfigServer
}

type ConfigServer struct {
	URL    string
	Routes []string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := Config{
		Address: ":3443",
		Certificate: &server.Certificate{
			CertFile: "./ipcmanview_data/cert.pem",
			KeyFile:  "./ipcmanview_data/key.pem",
		},
		Servers: []ConfigServer{
			{
				URL: "https://127.0.0.1:8443",
				Routes: []string{
					api.Route, api.Route + "/*",
					rpcserver.Route, rpcserver.Route + "/*",
				},
			},
			{
				URL:    "http://127.0.0.1:5173",
				Routes: []string{web.Route, web.Route + "*"},
			},
		},
	}

	start(ctx, cfg)
}

func start(ctx context.Context, cfg Config) {
	e := echo.New()
	e.Use(echoext.Logger())

	for _, server := range cfg.Servers {
		urL := must(url.Parse(server.URL))
		func(urL *url.URL) {
			for _, route := range server.Routes {
				e.Any(route, echo.WrapHandler(&httputil.ReverseProxy{
					Rewrite: func(r *httputil.ProxyRequest) {
						r.SetURL(urL)
						r.SetXForwarded()
					},
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				}))
			}
		}(urL)
	}

	err := server.NewHTTPServer(e, cfg.Address, cfg.Certificate).Serve(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal().Err(err).Send()
	}
}

func must[T any](d T, err error) T {
	if err != nil {
		panic(err)
	}
	return d
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel)
}
