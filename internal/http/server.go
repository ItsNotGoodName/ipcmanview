package http

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Server struct {
	e               *echo.Echo
	httpAddress     string
	httpsAddress    string
	https           bool
	cert            models.Certificate
	shutdownTimeout time.Duration
}

func NewServer(
	e *echo.Echo,
	address string,
) Server {
	return Server{
		e:               e,
		httpAddress:     address,
		shutdownTimeout: 3 * time.Second,
	}
}

func NewHTTPServer(
	e *echo.Echo,
	httpAddress string,
	httpsAddress string,
	https bool,
	cert models.Certificate,
) Server {
	return Server{
		e:               e,
		httpAddress:     httpAddress,
		httpsAddress:    httpsAddress,
		https:           https,
		cert:            cert,
		shutdownTimeout: 3 * time.Second,
	}
}

func (s Server) Serve(ctx context.Context) error {
	s.e.HideBanner = true
	s.e.HidePort = true
	log.Info().Str("address", s.httpAddress).Msg("Starting HTTP server")

	errC := make(chan error, 1)
	go func() {
		if s.https {
			errC <- s.e.StartTLS(s.httpsAddress, s.cert.CertFile, s.cert.KeyFile)
		} else {
			errC <- s.e.Start(s.httpAddress)
		}
	}()

	select {
	case err := <-errC:
		return errors.Join(suture.ErrTerminateSupervisorTree, err)
	case <-ctx.Done():
		log.Info().Msg("Gracefully shutting down HTTP server...")

		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := s.e.Shutdown(ctx); err != nil {
			log.Err(err).Msg("HTTP Server failed to shutdown gracefully")
			return err
		}

		return nil
	}
}
