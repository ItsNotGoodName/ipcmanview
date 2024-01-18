package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Page struct {
	DB repo.DB
}

// Profile implements rpc.Page.
func (p *Page) Profile(ctx context.Context, req *rpc.PageProfileReq) (*rpc.PageProfileResp, error) {
	sessionUser, err := useSessionUser(ctx)
	if err != nil {
		return nil, err
	}

	user, err := p.DB.GetUser(ctx, sessionUser.UserID)
	if err != nil {
		return nil, internalError(err)
	}

	return &rpc.PageProfileResp{
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
	}, nil
}
