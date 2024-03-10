package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewPublic() *Public {
	return &Public{}
}

type Public struct {
}

func (p *Public) GetConfig(context.Context, *emptypb.Empty) (*rpc.GetConfigResp, error) {
	cfg, err := system.GetConfig()
	if err != nil {
		return nil, err
	}

	return &rpc.GetConfigResp{
		SiteName:     cfg.SiteName,
		EnableSignUp: cfg.EnableSignUp,
	}, nil
}

func (p *Public) SignUp(ctx context.Context, req *rpc.SignUpReq) (*emptypb.Empty, error) {
	cfg, err := system.GetConfig()
	if err != nil {
		return nil, err
	}

	id, err := auth.CreateUser(ctx, cfg, auth.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errs, ok := core.AsFieldErrors(err); ok {
			return nil, newInvalidArgument(errs,
				keymap("email", "Email"),
				keymap("username", "Username"),
				keymap("password", "Password"),
			)
		}
		return nil, err
	}

	// TODO: remove this
	if err := auth.UpdateUserAdmin(core.WithSystemActor(ctx), id, true); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (*Public) ForgotPassword(context.Context, *rpc.ForgotPasswordReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}
