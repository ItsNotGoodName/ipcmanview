package sandbox

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/server"
	"github.com/ItsNotGoodName/ipcmanview/server/jwt"
	"github.com/ItsNotGoodName/ipcmanview/server/rpc"
	"github.com/ItsNotGoodName/ipcmanview/server/service"
	"github.com/ItsNotGoodName/ipcmanview/ui"
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
