package rpcserver

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (p *Page) GetHomePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetHomePageResp, error) {
	authSession := useAuthSession(ctx)

	dbDevices, err := p.db.ListDahuaDevicesForUser(ctx, repo.ListDahuaDevicesForUserParams{
		Admin:  authSession.Admin,
		UserID: core.Int64ToNullInt64(authSession.UserID),
	})
	if err != nil {
		return nil, check(err)
	}

	for _, lddfur := range dbDevices {
		fmt.Println(lddfur.ID, lddfur.Level)
	}

	return &rpc.GetHomePageResp{
		DeviceCount: int64(len(dbDevices)),
	}, nil
}

func (p *Page) GetProfilePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetProfilePageResp, error) {
	authSession := useAuthSession(ctx)

	user, err := p.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}

	dbSessions, err := p.db.ListUserSessionsForUserAndNotExpired(ctx, repo.ListUserSessionsForUserAndNotExpiredParams{
		UserID: authSession.UserID,
		Now:    types.NewTime(time.Now()),
	})
	if err != nil {
		return nil, check(err)
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

	return &rpc.GetProfilePageResp{
		Username:      user.Username,
		Email:         user.Email,
		Admin:         authSession.Admin,
		CreatedAtTime: timestamppb.New(user.CreatedAt.Time),
		UpdatedAtTime: timestamppb.New(user.UpdatedAt.Time),
		Sessions:      sessions,
	}, nil
}
