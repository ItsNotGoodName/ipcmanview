// main runs the development proxy.
package main

import (
	"context"
	"errors"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"

	"github.com/ItsNotGoodName/ipcmanview/internal/http"
	"github.com/ItsNotGoodName/ipcmanview/pkg/echoext"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	routes := []struct {
		URL   string
		Route string
	}{
		{
			URL:   "http://127.0.0.1:5174",
			Route: "/next/*",
		},
		{
			URL:   "http://127.0.0.1:5174",
			Route: "/next",
		},
		{
			URL:   "http://127.0.0.1:8080",
			Route: "/*",
		},
		{
			URL:   "http://127.0.0.1:8080",
			Route: "/",
		},
	}
	address := ":3000"

	e := echo.New()
	e.Use(echoext.Logger())

	for _, route := range routes {
		urL := must(url.Parse(route.URL))
		func(urL *url.URL) {
			e.Any(route.Route, echo.WrapHandler(&httputil.ReverseProxy{
				Rewrite: func(r *httputil.ProxyRequest) {
					r.SetURL(urL)
					r.SetXForwarded()
				},
			}))
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
