package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Page struct {
	DB repo.DB
}

func (p *Page) Profile(ctx context.Context, req *rpc.PageProfileReq) (*rpc.PageProfileResp, error) {
	authSession := useAuthSession(ctx)

	user, err := p.DB.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, NewError(err).Internal()
	}

	dbSessions, err := p.DB.ListUserSessionForUserAndNotExpired(ctx, repo.ListUserSessionForUserAndNotExpiredParams{
		UserID: authSession.UserID,
		Now:    types.NewTime(time.Now()),
	})
	if err != nil {
		return nil, NewError(err).Internal()
	}

	activeCutoff := time.Now().Add(-24 * time.Hour)
	sessions := make([]*rpc.Session, 0, len(dbSessions))
	for _, v := range dbSessions {
		sessions = append(sessions, &rpc.Session{
			Id:         v.ID,
			UserAgent:  v.UserAgent,
			Ip:         v.Ip,
			LastIp:     v.LastIp,
			LastUsedAt: timestamppb.New(v.LastUsedAt.Time),
			CreatedAt:  timestamppb.New(v.CreatedAt.Time),
			Active:     v.LastUsedAt.After(activeCutoff),
			Current:    v.Session == authSession.Session,
		})
	}

	return &rpc.PageProfileResp{
		Username:  user.Username,
		Email:     user.Email,
		Admin:     authSession.Admin,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
		Sessions:  sessions,
	}, nil
}
