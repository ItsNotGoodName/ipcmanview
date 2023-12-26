package http

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/internal/webnext"
	"github.com/ItsNotGoodName/ipcmanview/pkg/echoext"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	echoext.WithErrorLogging(e)

	// Middleware
	e.Use(echoext.LoggerWithConfig(echoext.LoggerConfig{
		Format: []string{
			"remote_ip",
			"host",
			"method",
			"user_agent",
			"status",
			"error",
			"latency_human",
		},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: echoext.RecoverLogErrorFunc,
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: web.AssetFS(),
	}))
	e.Group("/next", middleware.StaticWithConfig(middleware.StaticConfig{
		// Skipper: func(c echo.Context) bool {
		// 	// Prevent API 404's from being overwritten
		// 	return strings.HasPrefix(c.Request().RequestURI, "/api")
		// },
		Root:       "dist",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: webnext.DistFS(),
		IgnoreBase: true,
	}))

	return e
}
