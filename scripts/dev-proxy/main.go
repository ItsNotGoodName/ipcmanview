// main runs the development proxy.
package main

import (
	"context"
	"errors"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcserver"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/internal/webadmin"
	"github.com/ItsNotGoodName/ipcmanview/pkg/echoext"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	servers := []struct {
		URL    string
		Routes []string
	}{
		{
			URL: "http://127.0.0.1:8080",
			Routes: []string{
				webadmin.Route, webadmin.Route + "/*",
				api.Route, api.Route + "/*",
				rpcserver.Route, rpcserver.Route + "/*",
			},
		},
		{
			URL:    "http://127.0.0.1:5174",
			Routes: []string{web.Route, web.Route + "*"},
		},
	}
	address := ":3000"

	e := echo.New()
	e.Use(echoext.Logger())

	for _, server := range servers {
		urL := must(url.Parse(server.URL))
		func(urL *url.URL) {
			for _, route := range server.Routes {
				e.Any(route, echo.WrapHandler(&httputil.ReverseProxy{
					Rewrite: func(r *httputil.ProxyRequest) {
						r.SetURL(urL)
						r.SetXForwarded()
					},
				}))
			}
		}(urL)
	}

	err := http.NewServer(e, address).Serve(ctx)
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
