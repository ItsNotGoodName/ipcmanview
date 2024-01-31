package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewPublic(db repo.DB) *Public {
	return &Public{
		db: db,
	}
}

type Public struct {
	db repo.DB
}

func (p *Public) SignUp(ctx context.Context, req *rpc.SignUpReq) (*emptypb.Empty, error) {
	id, err := auth.CreateUser(ctx, p.db, auth.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errs, ok := asValidationErrors(err); ok {
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

	// TODO: remove this
	p.db.AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{UserID: id, CreatedAt: types.NewTime(time.Now())})

	return &emptypb.Empty{}, nil
}

func (*Public) ForgotPassword(context.Context, *rpc.ForgotPasswordReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}
