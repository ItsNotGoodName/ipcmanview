package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/go-playground/validator/v10"
)

func NewPublic(db repo.DB) *Public {
	return &Public{
		db: db,
	}
}

type Public struct {
	db repo.DB
}

func (p *Public) SignUp(ctx context.Context, req *rpc.SignUpReq) (*rpc.SignUpResp, error) {
	_, err := auth.CreateUser(ctx, p.db, auth.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return nil, NewError(err, "Failed to sign up.").Validation(errs, [][2]string{
				{"email", "Email"},
				{"username", "Username"},
				{"password", "Password"},
			})
		}

		if constraintErr, ok := asConstraintError(err); ok {
			return nil, NewError(err, "Failed to sign up.").Constraint(constraintErr, [][3]string{
				{"username", "users.username", "Name already taken."},
				{"email", "users.email", "Email already taken."},
			})
		}

		return nil, check(err)
	}

	return &rpc.SignUpResp{}, nil
}

func (*Public) ForgotPassword(context.Context, *rpc.ForgotPasswordReq) (*rpc.ForgotPasswordResp, error) {
	return nil, errNotImplemented
}
