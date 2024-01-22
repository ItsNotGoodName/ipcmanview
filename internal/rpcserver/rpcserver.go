package rpcserver

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
)

const Route = "/twirp"

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
			return ctx
		},
	})
}

// ---------- Middleware

func AuthSession() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			_, ok := auth.UseSession(ctx)
			if !ok {
				return ctx, twirp.Unauthenticated.Error("Invalid session or not signed in.")
			}
			return ctx, nil
		},
	})
}

func AdminAuthSession() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			authSession, ok := auth.UseSession(ctx)
			if !ok {
				return ctx, twirp.Unauthenticated.Error("Invalid session or not signed in.")
			}
			if !authSession.Admin {
				return ctx, twirp.PermissionDenied.Error("You are not an admin.")
			}

			return ctx, nil
		},
	})
}

func useAuthSession(ctx context.Context) models.AuthSession {
	u, ok := auth.UseSession(ctx)
	if !ok {
		panic("rpcserver.useAuthSession must be called after rpcserver.AuthSessionMiddleware")
	}
	return u
}

// ---------- Error

type Error struct {
	msg string
}

func NewError(err error, msg ...string) Error {
	if err != nil {
		log.Err(err).Str("package", "rpcserver").Send()
	}
	if len(msg) == 0 {
		return Error{msg: "Something went wrong."}
	}
	return Error{msg: msg[0]}
}

func (e Error) Field(field string, fieldErr error) twirp.Error {
	return twirp.InvalidArgument.Error(e.msg).WithMeta(field, fieldErr.Error())
}

func (e Error) Validation(errs validator.ValidationErrors, lookup [][2]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error(e.msg)
	for _, f := range errs {
		field := f.Field()
		for _, kv := range lookup {
			if kv[0] == field {
				twirpErr = twirpErr.WithMeta(kv[1], f.Error())
			}
		}
	}
	return twirpErr
}

func (e Error) Constraint(constraintErr sqlite.ConstraintError, lookup [][3]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error(e.msg)
	for _, kv := range lookup {
		if constraintErr.IsField(kv[0]) {
			twirpErr = twirpErr.WithMeta(kv[1], kv[2])
			break
		}
	}
	return twirpErr
}

func (w Error) Internal() twirp.Error {
	return twirp.InternalError(w.msg)
}

func (w Error) NotImplemented() twirp.Error {
	return twirp.InternalError("Not implemented.")
}
