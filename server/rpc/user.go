package rpc

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/db"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcgen"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
)

type UserService struct {
	db qes.Querier
}

var _ rpcgen.UserService = (*UserService)(nil)

func NewUserService(db qes.Querier) *UserService {
	return &UserService{
		db: db,
	}
}

// Me implements rpcgen.UserService.
func (a *UserService) Me(ctx context.Context) (*rpcgen.User, error) {
	// Get claim
	claim := auth.JWTClaimFromContext(ctx)

	// Get user
	user, err := db.User.Get(ctx, a.db, claim.UserID)
	if err != nil {
		return nil, handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	return convertUser(user), nil
}
