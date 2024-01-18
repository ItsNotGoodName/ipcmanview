package rpcserver

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/twitchtv/twirp"
)

var errNotImplemented = twirp.Internal.Error("Not implemented.")

type TwirpHandler interface {
	http.Handler
	PathPrefix() string
}

func Register(e *echo.Echo, t TwirpHandler, middleware ...echo.MiddlewareFunc) {
	e.Any(t.PathPrefix()+"*", echo.WrapHandler(t), middleware...)
}

func LoggerHooks() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		Error: func(ctx context.Context, err twirp.Error) context.Context {
			log.Err(err).Send()
			return ctx
		},
	})
}

func AuthHooks() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			user, ok := auth.GetSessionUser(ctx)
			if !ok || !user.Valid {
				return ctx, twirp.Unauthenticated.Error("Invalid session or not signed in.")
			}
			return ctx, nil
		},
	})
}

func useSessionUser(ctx context.Context) (auth.SessionUser, error) {
	u, ok := auth.GetSessionUser(ctx)
	if !ok {
		return auth.SessionUser{}, twirp.InternalError("Failed to get user by session.")
	}
	return u, nil
}

func validationError(errs validator.ValidationErrors, message string, lookup [][2]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error(message)
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

func constraintError(constraintErr sqlite.ConstraintError, message string, lookup [][3]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error(message)
	for _, kv := range lookup {
		if constraintErr.IsField(kv[0]) {
			twirpErr = twirpErr.WithMeta(kv[1], kv[2])
			break
		}
	}
	return twirpErr
}

func internalError(err error) twirp.Error {
	return twirp.InternalError(err.Error())
}
