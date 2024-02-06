package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewPublic(db sqlite.DB) *Public {
	return &Public{
		db: db,
	}
}

type Public struct {
	db sqlite.DB
}

func (p *Public) SignUp(ctx context.Context, req *rpc.SignUpReq) (*emptypb.Empty, error) {
	id, err := auth.CreateUser(ctx, p.db, auth.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	// TODO: remove this
	p.db.C().AuthUpsertAdmin(ctx, repo.AuthUpsertAdminParams{UserID: id, CreatedAt: types.NewTime(time.Now())})

	return &emptypb.Empty{}, nil
}

func (*Public) ForgotPassword(context.Context, *rpc.ForgotPasswordReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}
