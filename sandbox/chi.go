package sandbox

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmango/server"
	"github.com/ItsNotGoodName/ipcmango/server/service"
	"github.com/ItsNotGoodName/ipcmango/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func Chi(ctx context.Context) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	server.CORS(r)
	server.MountFS(r, ui.FS)

	r.Handle("/*", service.NewExampleServiceServer(ExampleServiceRPC{}))

	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-ctx.Done()

	err := s.Shutdown(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}
}

type ExampleServiceRPC struct{}

func (e ExampleServiceRPC) Message(ctx context.Context) (*service.Message, error) {
	return &service.Message{
		Body: "Hello from the server.",
		Time: time.Now(),
	}, nil
}
