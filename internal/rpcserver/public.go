package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/config"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewPublic(configProvider config.Provider, db sqlite.DB) *Public {
	return &Public{
		configProvider: configProvider,
		db:             db,
	}
}

type Public struct {
	configProvider config.Provider
	db             sqlite.DB
}

func (p *Public) GetConfig(context.Context, *emptypb.Empty) (*rpc.GetConfigResp, error) {
	cfg, err := p.configProvider.GetConfig()
	if err != nil {
		return nil, err
	}

	return &rpc.GetConfigResp{
		SiteName:     cfg.SiteName,
		EnableSignUp: cfg.EnableSignUp,
	}, nil
}

func (p *Public) SignUp(ctx context.Context, req *rpc.SignUpReq) (*emptypb.Empty, error) {
	cfg, err := p.configProvider.GetConfig()
	if err != nil {
		return nil, err
	}

	id, err := auth.CreateUser(ctx, cfg, p.db, auth.CreateUserParams{
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
