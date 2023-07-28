package sandbox

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmango/server/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Chi() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Handle("/*", service.NewExampleServiceServer(ExampleServiceRPC{}))

	http.ListenAndServe(":3000", r)
}

type ExampleServiceRPC struct{}

func (e ExampleServiceRPC) Message(ctx context.Context) (*service.Message, error) {
	return &service.Message{
		Body: "Hello from the server.",
		Time: time.Now(),
	}, nil
}
