//go:build dev

package http

import (
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
)

func useWebNext(e *echo.Echo) {
	u, err := url.Parse("http://127.0.0.1:5174")
	if err != nil {
		panic(err)
	}

	h := echo.WrapHandler(&httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(u)
		},
	})

	e.Group("/next").Any("*", h)
}
