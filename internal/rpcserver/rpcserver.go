package rpcserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TwirpHandler interface {
	http.Handler
	PathPrefix() string
}

func Register(e *echo.Echo, t TwirpHandler) {
	e.Any(t.PathPrefix()+"*", echo.WrapHandler(t))
}
