package sandbox

import (
	"context"

	"github.com/ItsNotGoodName/ipcmango/server"
	"github.com/ItsNotGoodName/ipcmango/server/jwt"
	"github.com/ItsNotGoodName/ipcmango/server/rpc"
	"github.com/ItsNotGoodName/ipcmango/server/service"
	"github.com/ItsNotGoodName/ipcmango/ui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Chi(ctx context.Context, shutdown context.CancelFunc, pool *pgxpool.Pool) server.HTTP {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	server.CORS(r)
	server.MountFS(r, ui.FS)

	userService := rpc.NewUserService(pool)
	r.Handle("/rpc/AuthService/*", service.NewAuthServiceServer(userService))

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwt.TokenAuth))
		r.Use(jwt.Authenticator)

		r.Handle("/rpc/UserService/*", service.NewUserServiceServer(userService))
	})

	return server.NewHTTP(r, ":8080", shutdown)
}
