package sandbox

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmango/server"
	"github.com/ItsNotGoodName/ipcmango/server/service"
	"github.com/ItsNotGoodName/ipcmango/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Chi(ctx context.Context, shutdown context.CancelFunc) server.HTTP {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	server.CORS(r)
	server.MountFS(r, ui.FS)

	r.Handle("/*", service.NewExampleServiceServer(ExampleServiceRPC{}))

	return server.NewHTTP(r, ":8080", shutdown)
}

type ExampleServiceRPC struct{}

func (e ExampleServiceRPC) Message(ctx context.Context) (*service.Message, error) {
	return &service.Message{
		Body: "Hello from the server.",
		Time: time.Now(),
	}, nil
}
