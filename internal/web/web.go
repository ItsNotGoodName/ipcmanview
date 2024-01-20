// webnext contains the SPA.
package web

import (
	"embed"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate pnpm install
//go:generate pnpm run build

//go:embed dist
var dist embed.FS

const Route = "/"

func FS(skipPrefix ...string) echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			// Prevent API 404's from being overwritten
			for _, prefix := range skipPrefix {
				if strings.HasPrefix(c.Request().RequestURI, prefix) {
					return true
				}
			}
			return false
		},
		Root:       "dist",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(dist),
		IgnoreBase: true,
	})
}
