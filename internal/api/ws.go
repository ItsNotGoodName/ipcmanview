package api

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/apiws"
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
)

type WSData struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type WSEvent struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

func (s Server) WS(c echo.Context) error {
	w := c.Response()
	r := c.Request()
	ctx := r.Context()

	conn, err := apiws.Upgrade(w, r)
	if err != nil {
		return err
	}

	WS(ctx, conn, s.pub)

	return nil
}

func WS(ctx context.Context, conn *websocket.Conn, pub *pubsub.Pub) {
	actor := core.UseActor(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := apiws.Logger(conn)

	sub, eventC, err := pub.
		Subscribe().
		Middleware(dahua.PubSubMiddleware(ctx)).
		Channel(ctx, 1)
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer sub.Close()

	// Visitors
	final := apiws.NewOnceVisitor()
	buffer := apiws.NewBufferVisitor(100)
	visitors := apiws.NewVisitors(final, buffer)

	// IO
	sig := apiws.NewSignal()
	writerC := apiws.Writer(ctx, conn, log, sig)
	readC := apiws.Reader(ctx, conn, log)

	for {
		apiws.Check(visitors, sig)

		select {
		case <-ctx.Done():
			return
		case data, ok := <-readC:
			// Read
			if !ok {
				return
			}

			log.Error().Bytes("data", data).Msg("WebSocket client is not supposed to send data...")
		case writeC, ok := <-writerC:
			// Write
			if !ok {
				return
			}

			if err := apiws.Flush(ctx, visitors, writeC); err != nil {
				log.Err(err).Msg("Failed to flush")
				return
			}

			if final.Done {
				return
			}
		case event, ok := <-eventC:
			// Pub sub
			if !ok {
				return
			}

			var isFinal bool
			var payload WSData
			switch event := event.(type) {
			case bus.DahuaEmailCreated:
				payload = WSData{
					Type: "event",
					Data: WSEvent{
						Action: "dahua-email:created",
						Data:   event.MessageID,
					},
				}
			case bus.UserSecurityUpdated:
				payload = WSData{
					Type: "event",
					Data: WSEvent{
						Action: "user-security:updated",
						Data:   event.UserID,
					},
				}

				if event.UserID == actor.UserID {
					isFinal = true
				}
			case bus.DahuaFileCreated:
				payload = WSData{
					Type: "event",
					Data: WSEvent{
						Action: "dahua-scan-file:created",
						Data:   event.Count,
					},
				}
			case bus.DahuaFileCursorUpdated:
				payload = WSData{
					Type: "event",
					Data: WSEvent{
						Action: "dahua-file-cursor:updated",
						Data:   event.Cursor,
					},
				}
			case bus.DahuaEvent:
				if event.EventRule.IgnoreLive {
					continue
				}

				payload = WSData{
					Type: "dahua-event",
					Data: dahua.NewDahuaEvent(event.Event),
				}
			}

			b, err := json.Marshal(payload)
			if err != nil {
				log.Err(err).Send()
				return
			}

			if isFinal {
				final.Set(b)
			} else if !buffer.Push(b) {
				return
			}
		}
	}
}
