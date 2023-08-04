package rpc

import (
	"context"
	"strings"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/server/jwt"
	"github.com/ItsNotGoodName/ipcmango/server/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool: pool,
	}
}

func (u UserService) context(ctx context.Context) (*pgxpool.Conn, func(), error) {
	conn, err := u.pool.Acquire(ctx)
	if err != nil {
		return nil, nil, service.ErrorWithCause(service.ErrWebrpcServerPanic, err)
	}
	return conn, conn.Release, nil
}

func (u UserService) Login(ctx context.Context, usernameOrEmail string, password string) (*service.User, string, error) {
	conn, release, err := u.context(ctx)
	if err != nil {
		return nil, "", err
	}
	defer release()

	user, err := db.User.GetByUsernameOrEmail(ctx, conn, strings.ToLower(usernameOrEmail))
	if err != nil {
		return nil, "", service.ErrWebrpcBadRequest
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, "", service.ErrWebrpcBadRequest
	}

	return newUser(user), jwt.EncodeUserID(user.ID), nil
}

func (u UserService) Me(ctx context.Context) (*service.User, error) {
	conn, release, err := u.context(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

	id := jwt.DecodeUserID(ctx)
	user, err := db.User.Get(ctx, conn, id)
	if err != nil {
		return nil, service.ErrorWithCause(service.ErrWebrpcBadResponse, err)
	}

	return newUser(user), nil
}

func (u UserService) Register(ctx context.Context, r *service.UserRegister) error {
	conn, release, err := u.context(ctx)
	if err != nil {
		return err
	}
	defer release()

	user, err := core.NewUser(core.UserCreate{
		Email:           r.Email,
		Username:        r.Username,
		Password:        r.Password,
		PasswordConfirm: r.PasswordConfirm,
	})
	if err != nil {
		return service.ErrorWithCause(service.ErrWebrpcBadRequest, err)
	}

	user, err = db.User.Create(ctx, conn, user)
	if err != nil {
		return service.ErrorWithCause(service.ErrWebrpcBadRequest, err)
	}

	return nil
}
