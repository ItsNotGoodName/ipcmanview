package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type HTTP struct {
	r        chi.Router
	addr     string
	shutdown context.CancelFunc
}

func NewHTTP(r chi.Router, addr string, shutdown context.CancelFunc) HTTP {
	return HTTP{
		r:        r,
		addr:     addr,
		shutdown: shutdown,
	}
}

func (r HTTP) Serve(ctx context.Context) error {
	server := &http.Server{Addr: r.addr, Handler: r.r}

	go func() {
		<-ctx.Done()

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				server.Close()
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Err(err).Msg("Failed to shutdown HTTP server")
		}
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Err(err).Msg("Failed to start HTTP server")
		r.shutdown()
	}

	return nil
}
