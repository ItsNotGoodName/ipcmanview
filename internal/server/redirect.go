package server

import (
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/labstack/echo/v4"
)

func NewRedirect(httpsPort string) *echo.Echo {
	e := echo.New()

	e.Any("*", func(c echo.Context) error {
		r := c.Request()

		host, _ := core.SplitAddress(r.Host)

		http.Redirect(c.Response(), r, "https://"+host+":"+httpsPort+r.RequestURI, http.StatusMovedPermanently)
		return nil
	})

	return e
}
