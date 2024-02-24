package rpcserver

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
)

const Route = "/twirp"

var errNotImplemented = twirp.InternalError("Not implemented.")

// ---------- Server

func NewServer(e *echo.Echo) Server {
	return Server{e: e}
}

type Server struct {
	e *echo.Echo
}

type TwirpHandler interface {
	http.Handler
	PathPrefix() string
}

func (s Server) Register(t TwirpHandler, middleware ...echo.MiddlewareFunc) Server {
	s.e.Any(t.PathPrefix()+"*", echo.WrapHandler(t), middleware...)
	return s
}

func Logger() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		Error: func(ctx context.Context, err twirp.Error) context.Context {
			if err != nil {
				log.Err(err).Str("package", "rpcserver").Send()
			}
			return ctx
		},
	})
}

// ---------- Middleware

// RequireAuthSession allows only valid sessions.
func RequireAuthSession() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			session, ok := auth.UseSession(ctx)
			if !ok {
				return ctx, twirp.Unauthenticated.Error("Invalid session or not signed in.")
			}
			if session.Disabled {
				return ctx, twirp.Unauthenticated.Error("Account disabled.")
			}
			return ctx, nil
		},
	})
}

// RequireAuthSession allows only valid admin sessions.
func RequireAdminAuthSession() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			session, ok := auth.UseSession(ctx)
			if !ok {
				return ctx, twirp.Unauthenticated.Error("Invalid session or not signed in.")
			}
			if session.Disabled {
				return ctx, twirp.Unauthenticated.Error("Account disabled.")
			}
			if !session.Admin {
				return ctx, twirp.PermissionDenied.Error("You are not an admin.")
			}
			return ctx, nil
		},
	})
}

func useAuthSession(ctx context.Context) auth.Session {
	u, ok := auth.UseSession(ctx)
	if !ok {
		panic("rpcserver.useAuthSession must be called after rpcserver.RequireAuthSession or rpcserver.RequireAdminAuthSession")
	}
	return u
}

// ---------- Error

func keymap(external, internal string, message ...string) [3]string {
	if len(message) > 0 {
		return [3]string{external, internal, message[0]}
	} else {
		return [3]string{external, internal, ""}
	}
}

func newInvalidArgument(errs core.FieldErrors, keymap ...[3]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error("Invalid argument.")
	for _, f := range errs {
		for _, km := range keymap {
			if km[1] == f.Field {
				if km[2] == "" {
					twirpErr = twirpErr.WithMeta(km[0], f.Message())
				} else {
					twirpErr = twirpErr.WithMeta(km[0], km[2])
				}
			}
		}
	}
	return twirpErr
}
