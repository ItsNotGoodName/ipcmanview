package api

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/apiws"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/gorilla/websocket"
)

type WSData struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type WSEvent struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

func WS(ctx context.Context, conn *websocket.Conn, db sqlite.DB, pub *pubsub.Pub) {
	actor := core.UseActor(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := apiws.Logger(conn)

	sub, eventC, err := pub.
		Subscribe(event.Event{}, event.DahuaEvent{}).
		Middleware(dahua.PubSubMiddleware(ctx, db)).
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
	writerC := apiws.Writer(ctx, cancel, conn, log, sig)
	readC := apiws.Reader(ctx, cancel, conn, log)

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
		case writeC := <-writerC:
			// Write
			if err := apiws.Flush(ctx, visitors, writeC); err != nil {
				log.Err(err).Msg("Failed to flush")
				return
			}

			if final.Done {
				return
			}
		case evt, ok := <-eventC:
			// Pub sub
			if !ok {
				return
			}

			var isFinal bool
			var payload WSData
			switch evt := evt.(type) {
			case event.Event:
				payload = WSData{
					Type: "event",
					Data: WSEvent{
						Action: string(evt.Event.Action),
						Data:   evt.Event.Data.RawMessage,
					},
				}

				if evt.Event.Action == event.ActionUserSecurityUpdated && event.DataAsInt64(evt.Event) == actor.UserID {
					isFinal = true
				}
			case event.DahuaEvent:
				if evt.EventRule.IgnoreLive {
					continue
				}

				payload = WSData{
					Type: "dahua-event",
					Data: dahua.NewDahuaEvent(evt.Event),
				}
			}

			b, err := json.Marshal(payload)
			if err != nil {
				log.Err(err).Send()
				return
			}

			if !buffer.Push(b) {
				return
			}

			if isFinal {
				final.Set(b)
			}
		}
	}
}
