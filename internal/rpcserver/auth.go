package rpcserver

import (
	"context"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/go-playground/validator/v10"
	"github.com/twitchtv/twirp"
)

type Auth struct {
	DB repo.DB
}

// ResetPassword implements rpc.Auth.
func (*Auth) ResetPassword(context.Context, *rpc.AuthResetPasswordReq) (*rpc.AuthResetPasswordResp, error) {
	return nil, errNotImplemented
}

// SignIn implements rpc.Auth.
func (a *Auth) SignIn(ctx context.Context, req *rpc.AuthSignInReq) (*rpc.AuthSignInResp, error) {
	user, err := a.DB.GetUserByUsernameOrEmail(ctx, strings.ToLower(strings.TrimSpace(req.UsernameOrEmail)))
	if err != nil {
		return nil, twirp.Internal.Error("Incorrect credentials.")
	}

	if err := auth.CheckUserPassword(user.Password, req.Password); err != nil {
		return nil, twirp.Internal.Error("Incorrect credentials.")
	}

	session, err := auth.CreateSesssion(ctx, a.DB, user.ID, auth.DefaultSessionDuration)
	if err != nil {
		return nil, twirp.Internal.Error("Something went wrong.")
	}

	return &rpc.AuthSignInResp{
		Token: session,
	}, nil
}

// SignUp implements rpc.Auth.
func (a *Auth) SignUp(ctx context.Context, req *rpc.AuthSignUpReq) (*rpc.AuthSignUpResp, error) {
	_, err := auth.CreateUser(ctx, a.DB, models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return nil, validationError(errs, "Failed to sign up.", [][2]string{
				{"Email", "email"},
				{"Username", "username"},
				{"Password", "password"},
			})
		}

		if constraintErr, ok := sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE); ok {
			return nil, constraintError(constraintErr, "User already exists.", [][3]string{
				{"users.username", "username", "Name already taken."},
				{"users.email", "email", "Email already taken."},
			})
		}

		return nil, err
	}

	return &rpc.AuthSignUpResp{}, nil
}
