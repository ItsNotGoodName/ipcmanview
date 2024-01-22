package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/go-playground/validator/v10"
)

func NewAuth(db repo.DB) *Auth {
	return &Auth{
		db: db,
	}
}

type Auth struct {
	db repo.DB
}

func (a *Auth) SignUp(ctx context.Context, req *rpc.AuthSignUpReq) (*rpc.AuthSignUpResp, error) {
	_, err := auth.CreateUser(ctx, a.db, models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return nil, NewError(err, "Failed to sign up.").Validation(errs, [][2]string{
				{"Email", "email"},
				{"Username", "username"},
				{"Password", "password"},
			})
		}

		if constraintErr, ok := sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE); ok {
			return nil, NewError(err, "Failed to sign up.").Constraint(constraintErr, [][3]string{
				{"users.username", "username", "Name already taken."},
				{"users.email", "email", "Email already taken."},
			})
		}

		return nil, NewError(err).Internal()
	}

	return &rpc.AuthSignUpResp{}, nil
}

func (*Auth) Forgot(context.Context, *rpc.AuthForgotReq) (*rpc.AuthForgotResp, error) {
	return nil, NewError(nil).NotImplemented()
}
