package server

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcgen"
	"github.com/ItsNotGoodName/ipcmanview/server/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Router(
	jwt auth.JWTAuth,
	h api.Handler,
	authService rpcgen.AuthService,
	userService rpcgen.UserService,
	dahauService rpcgen.DahuaService,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Handle("/rpc/AuthService/*", rpcgen.NewAuthServiceServer(authService))
	r.Group(func(r chi.Router) {
		r.Use(jwt.Verifier)
		r.Use(auth.JWTAuthenticator)

		r.Handle("/rpc/UserService/*", rpcgen.NewUserServiceServer(userService))
		r.Handle("/rpc/DahuaService/*", rpcgen.NewDahuaServiceServer(dahauService))
	})

	r.Group(func(r chi.Router) {
		r.Use(jwt.Verifier)
		r.Use(auth.JWTAuthenticator)

		r.Get("/v1/dahua/cameras/{id}/snapshot", h.WithID(api.Snapshot))
	})

	return r
}
