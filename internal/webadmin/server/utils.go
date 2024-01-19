package server

import (
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/pkg/htmx"
	"github.com/labstack/echo/v4"
)

func toast(message string) htmx.Event {
	return htmx.NewEvent("toast", message)
}

func toastError(message string) htmx.Event {
	return htmx.NewEvent("toast-error", message)
}

// isHTMX checks if request is an htmx request but not a boosted htmx request.
func isHTMX(c echo.Context) bool {
	return htmx.GetRequest(c.Request()) && !htmx.GetBoosted(c.Request())
}

func useEventStream(c echo.Context) *echo.Response {
	w := c.Response()
	w.Header().Set(echo.HeaderContentType, "text/event-stream")
	w.Header().Set(echo.HeaderCacheControl, "no-cache")
	w.Header().Set(echo.HeaderConnection, "keep-alive")
	return w
}

func sendEventStream(w *echo.Response, data []byte) error {
	_, err := w.Write(data)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func formatEventStream(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}
