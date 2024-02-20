package api

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/apiws"
	"github.com/gorilla/websocket"
)

func WS(ctx context.Context, conn *websocket.Conn) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log := apiws.Logger(conn)

	// Visitors
	buffer := apiws.NewBufferVisitor(10)
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

			apiws.Check(visitors, sig)
		}
	}
}
