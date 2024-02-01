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
		http.Redirect(c.Response(), r, "https://"+core.SplitAddress(r.Host)[0]+":"+httpsPort+r.RequestURI, http.StatusMovedPermanently)
		return nil
	})

	return e
}
