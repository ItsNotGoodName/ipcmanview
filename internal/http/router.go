package http

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
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
	useWebNext(e)

	return e
}
