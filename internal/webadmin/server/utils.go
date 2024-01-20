package server

import (
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/pkg/htmx"
	"github.com/labstack/echo/v4"
)

// isHTMX checks if request is an htmx request but not a boosted htmx request.
func isHTMX(c echo.Context) bool {
	return htmx.GetRequest(c.Request()) && !htmx.GetBoosted(c.Request())
}

// ---------- Toast

func toast(message string) htmx.Event {
	return htmx.NewEvent("toast", message)
}

func toastError(message string) htmx.Event {
	return htmx.NewEvent("toast-error", message)
}

// ---------- EventStream

func newEventStream(c echo.Context) *echo.Response {
	w := c.Response()
	h := w.Header()
	h.Set(echo.HeaderContentType, "text/event-stream")
	h.Set(echo.HeaderCacheControl, "no-cache")
	h.Set(echo.HeaderConnection, "keep-alive")
	return w
}

func writeEventStream(w *echo.Response, event, message string) error {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(message, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	data := []byte(eventPayload + "\n")

	_, err := w.Write(data)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
