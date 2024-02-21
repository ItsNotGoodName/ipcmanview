package api

import (
	"context"
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/apiws"
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
	buffer := apiws.NewBufferVisitor(100)
	visitors := apiws.NewVisitors(buffer)

	// IO
	sig := apiws.NewSignal()
	writerC := apiws.Writer(ctx, cancel, conn, log, sig)
	readC := apiws.Reader(ctx, cancel, conn, log)

	for {
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
		case evt, ok := <-eventC:
			// Pub sub
			if !ok {
				return
			}

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
		}

		apiws.Check(visitors, sig)
	}
}
