package webserver

import (
	"strings"

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
