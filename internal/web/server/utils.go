package webserver

import (
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

type Data map[string]any

var decoder = schema.NewDecoder()

func parseForm(c echo.Context, form any) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	if err := decoder.Decode(form, c.Request().PostForm); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return nil
}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}

func useDahuaCamera(c echo.Context, db *sqlc.Queries) (sqlc.DahuaCamera, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return sqlc.DahuaCamera{}, echo.ErrBadRequest.WithInternal(err)
	}

	camera, err := db.GetDahuaCamera(c.Request().Context(), id)
	if err != nil {
		return sqlc.DahuaCamera{}, err
	}

	return camera, nil
}
