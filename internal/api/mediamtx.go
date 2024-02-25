package api

import (
	"net/http"
	"net/http/httputil"

	echo "github.com/labstack/echo/v4"
)

func (s *Server) Mediamtx(prefix string) echo.HandlerFunc {
	return echo.WrapHandler(http.StripPrefix(prefix, httputil.NewSingleHostReverseProxy(s.mediamtxURL)))
}
