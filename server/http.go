package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type HTTP struct {
	r       chi.Router
	address string
}

func NewHTTP(r chi.Router, addr string) HTTP {
	return HTTP{
		r:       r,
		address: addr,
	}
}

func (r HTTP) Serve(ctx context.Context) error {
	server := &http.Server{Addr: r.address, Handler: r.r}

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

	log.Info().Str("address", r.address).Msg("Starting HTTP server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
