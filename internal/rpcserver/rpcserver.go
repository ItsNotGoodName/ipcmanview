package rpcserver

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/labstack/echo/v4"
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
			return ctx
		},
	})
}

// ---------- Middleware

func AuthSession() twirp.ServerOption {
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

func AdminAuthSession() twirp.ServerOption {
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
		panic("rpcserver.useAuthSession must be called after rpcserver.AuthSessionMiddleware")
	}
	return u
}

// func useDahuaPermissions(ctx context.Context, db sqlite.DB) (models.DahuaDevicePermissions, error) {
// 	session := useAuthSession(ctx)
// 	permissions, err := db.DahuaListDahuaDevicePermissions(ctx, session.UserID)
// 	return permissions, err
// }

// ---------- Error

// func asValidationErrors(err error) (validator.ValidationErrors, bool) {
// 	errs, ok := err.(validator.ValidationErrors)
// 	return errs, ok
// }
//
// func asConstraintError(err error) (sqlite.ConstraintError, bool) {
// 	return sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE)
// }
//
// type Error struct {
// 	msg string
// }
//
// func NewError(err error, msg string) Error {
// 	if err != nil {
// 		log.Err(err).Str("package", "rpcserver").Send()
// 	}
// 	return Error{msg: msg}
// }
//
// func (e Error) Field(field string) twirp.Error {
// 	return twirp.InvalidArgument.Error(e.msg).WithMeta(field, e.msg)
// }
//
// func (e Error) Validation(errs validator.ValidationErrors, lookup [][2]string) twirp.Error {
// 	twirpErr := twirp.InvalidArgument.Error(e.msg)
// 	for _, f := range errs {
// 		field := f.Field()
// 		for _, kv := range lookup {
// 			if kv[1] == field {
// 				twirpErr = twirpErr.WithMeta(kv[0], f.Translate(core.Translator))
// 			}
// 		}
// 	}
// 	return twirpErr
// }
//
// func (e Error) Constraint(constraintErr sqlite.ConstraintError, lookup [][3]string) twirp.Error {
// 	twirpErr := twirp.InvalidArgument.Error(e.msg)
// 	for _, kv := range lookup {
// 		if constraintErr.IsField(kv[1]) {
// 			twirpErr = twirpErr.WithMeta(kv[0], kv[2])
// 			break
// 		}
// 	}
// 	return twirpErr
// }
//
// func (w Error) Internal() twirp.Error {
// 	return twirp.InternalError(w.msg)
// }
//
// func (w Error) NotFound() twirp.Error {
// 	return twirp.NotFoundError(w.msg)
// }
