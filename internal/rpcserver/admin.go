package rpcserver

import (
	"context"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
)

func NewAdmin(db repo.DB) *Admin {
	return &Admin{
		db: db,
	}
}

type Admin struct {
	db repo.DB
}

func (a *Admin) ListGroups(ctx context.Context, req *rpc.ListGroupsReq) (*rpc.ListGroupsResp, error) {
	return nil, NewError(nil).NotImplemented()
}
