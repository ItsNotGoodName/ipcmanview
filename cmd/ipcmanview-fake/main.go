// ipcmanview-fake is used to develop the UI without the server.
package main

import (
	"context"
	"errors"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/server"
	"github.com/ItsNotGoodName/ipcmanview/server/jwt"
	"github.com/ItsNotGoodName/ipcmanview/server/service"
	"github.com/ItsNotGoodName/ipcmanview/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx, shutdown := interrupt.Context()
	defer shutdown()

	svc := NewService()

	// Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	server.CORS(r)
	server.MountFS(r, ui.FS)

	r.Handle("/rpc/AuthService/*", service.NewAuthServiceServer(svc))
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt.TokenAuth))
		r.Use(jwt.Authenticator)

		r.Handle("/rpc/UserService/*", service.NewUserServiceServer(svc))
		r.Handle("/rpc/DahuaService/*", service.NewDahuaServiceServer(svc))
	})

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// HTTP
	super.Add(server.NewHTTP(r, ":8080"))

	if err := super.Serve(ctx); !errors.Is(err, context.Canceled) {
		log.Err(err).Msg("Failed to start root supervisor")
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
