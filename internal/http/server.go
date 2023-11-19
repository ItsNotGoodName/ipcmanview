package http

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e               *echo.Echo
	address         string
	shutdownTimeout time.Duration
}

func NewServer(echo *echo.Echo, address string) Server {
	return Server{
		e:               echo,
		address:         address,
		shutdownTimeout: 3 * time.Second,
	}
}

func (s Server) Serve(ctx context.Context) error {
	errC := make(chan error, 1)
	go func() {
		errC <- s.e.Start(s.address)
	}()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		return s.e.Shutdown(ctx)
	}
}
