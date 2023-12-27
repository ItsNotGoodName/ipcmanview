//go:build !dev

package http

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/webnext"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func useWebNext(e *echo.Echo) {
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
}
