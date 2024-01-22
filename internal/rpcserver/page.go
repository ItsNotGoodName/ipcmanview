package rpcserver

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPage(db repo.DB) *Page {
	return &Page{
		db: db,
	}
}

type Page struct {
	db repo.DB
}

func (p *Page) Home(ctx context.Context, req *rpc.PageHomeReq) (*rpc.PageHomeResp, error) {
	authSession := useAuthSession(ctx)

	dbDevices, err := p.db.ListDahuaDeviceForUser(ctx, repo.ListDahuaDeviceForUserParams{
		Admin:  authSession.Admin,
		UserID: core.Int64ToNullInt64(authSession.UserID),
	})
	if err != nil {
		return nil, NewError(err).Internal()
	}

	for _, lddfur := range dbDevices {
		fmt.Println(lddfur.ID, lddfur.Level)
	}

	return &rpc.PageHomeResp{
		DeviceCount: int64(len(dbDevices)),
	}, nil
}

func (p *Page) Profile(ctx context.Context, req *rpc.PageProfileReq) (*rpc.PageProfileResp, error) {
	authSession := useAuthSession(ctx)

	user, err := p.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, NewError(err).Internal()
	}

	dbSessions, err := p.db.ListUserSessionForUserAndNotExpired(ctx, repo.ListUserSessionForUserAndNotExpiredParams{
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
			Id:             v.ID,
			UserAgent:      v.UserAgent,
			Ip:             v.Ip,
			LastIp:         v.LastIp,
			LastUsedAtTime: timestamppb.New(v.LastUsedAt.Time),
			CreatedAtTime:  timestamppb.New(v.CreatedAt.Time),
			Active:         v.LastUsedAt.After(activeCutoff),
			Current:        v.Session == authSession.Session,
		})
	}

	return &rpc.PageProfileResp{
		Username:      user.Username,
		Email:         user.Email,
		Admin:         authSession.Admin,
		CreatedAtTime: timestamppb.New(user.CreatedAt.Time),
		UpdatedAtTime: timestamppb.New(user.UpdatedAt.Time),
		Sessions:      sessions,
	}, nil
}
