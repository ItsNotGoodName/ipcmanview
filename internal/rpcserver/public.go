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
		return nil, checkCreateUpdateUser(err, "Failed to sign up.")
	}

	// TODO: remove this
	p.db.AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{UserID: id, CreatedAt: types.NewTime(time.Now())})

	return &emptypb.Empty{}, nil
}

func (*Public) ForgotPassword(context.Context, *rpc.ForgotPasswordReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}
