package rpcserver

import (
	"context"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/go-playground/validator/v10"
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
			return ctx
		},
	})
}

// ---------- Middleware

func AuthSession() twirp.ServerOption {
	return twirp.WithServerHooks(&twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			authSession, ok := auth.UseSession(ctx)
			if !ok {
				return ctx, twirp.Unauthenticated.Error("Invalid session or not signed in.")
			}
			if authSession.Disabled {
				return ctx, twirp.Unauthenticated.Error("Account disabled.")
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
			if authSession.Disabled {
				return ctx, twirp.Unauthenticated.Error("Account disabled.")
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

func asValidationErrors(err error) (validator.ValidationErrors, bool) {
	errs, ok := err.(validator.ValidationErrors)
	return errs, ok
}

func asConstraintError(err error) (sqlite.ConstraintError, bool) {
	return sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE)
}

func check(err error) twirp.Error {
	if repo.IsNotFound(err) {
		return NewError(err, "Not found.").NotFound()
	}
	return NewError(err, "Something went wrong.").Internal()
}

type Error struct {
	msg string
}

func NewError(err error, msg string) Error {
	if err != nil {
		log.Err(err).Str("package", "rpcserver").Send()
	}
	return Error{msg: msg}
}

func (e Error) Field(field string) twirp.Error {
	return twirp.InvalidArgument.Error(e.msg).WithMeta(field, e.msg)
}

func (e Error) Validation(errs validator.ValidationErrors, lookup [][2]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error(e.msg)
	for _, f := range errs {
		field := f.Field()
		for _, kv := range lookup {
			if kv[1] == field {
				twirpErr = twirpErr.WithMeta(kv[0], f.Error())
			}
		}
	}
	return twirpErr
}

func (e Error) Constraint(constraintErr sqlite.ConstraintError, lookup [][3]string) twirp.Error {
	twirpErr := twirp.InvalidArgument.Error(e.msg)
	for _, kv := range lookup {
		if constraintErr.IsField(kv[1]) {
			twirpErr = twirpErr.WithMeta(kv[0], kv[2])
			break
		}
	}
	return twirpErr
}

func (w Error) Internal() twirp.Error {
	return twirp.InternalError(w.msg)
}

func (w Error) NotFound() twirp.Error {
	return twirp.NotFoundError(w.msg)
}

// ---------- Convert/Parse

func parsePagePagination(v *rpc.PagePagination) pagination.Page {
	var (
		page    int
		perPage int
	)
	if v != nil {
		page = int(v.Page)
		perPage = int(v.PerPage)
	}

	if page < 1 {
		page = 1
	}
	if v.PerPage < 1 || v.PerPage > 100 {
		perPage = 10
	}

	return pagination.Page{
		Page:    page,
		PerPage: perPage,
	}
}

func convertPagePaginationResult(v pagination.PageResult) *rpc.PagePaginationResult {
	return &rpc.PagePaginationResult{
		Page:         int32(v.Page),
		PerPage:      int32(v.PerPage),
		TotalPages:   int32(v.TotalPages),
		TotalItems:   int64(v.TotalItems),
		SeenItems:    int64(v.Seen()),
		PreviousPage: int32(v.Previous()),
		NextPage:     int32(v.Next()),
	}
}

func convertOrderToSQL(sql string, o rpc.Order) string {
	switch o {
	case rpc.Order_DESC:
		return sql + " DESC"
	case rpc.Order_ASC:
		return sql + " ASC"
	default:
		return sql
	}
}
