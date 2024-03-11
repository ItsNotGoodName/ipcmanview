package dahuasmtp

import (
	"context"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

type Option func(s *smtp.Server)

func MaxMessageBytes(maxMessageBytes int64) Option {
	return func(s *smtp.Server) {
		s.MaxMessageBytes = maxMessageBytes
	}
}

type Server struct {
	server *smtp.Server
}

func NewServer(backend smtp.Backend, address string, options ...Option) Server {
	server := smtp.NewServer(backend)

	server.Addr = address
	server.Domain = "localhost"
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 25 * 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	for _, option := range options {
		option(server)
	}

	enableMechLogin(backend, server)

	return Server{
		server: server,
	}
}

func (Server) String() string {
	return "dahuasmtp.Server"
}

func (s Server) Serve(ctx context.Context) error {
	log.Info().Str("address", s.server.Addr).Msg("Starting SMTP server")

	errC := make(chan error, 1)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			errC <- err
		}
	}()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
	}

	log.Info().Msg("Gracefully shutting down SMTP server...")

	s.server.Close()

	return nil
}
