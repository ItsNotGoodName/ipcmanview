package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewUser(db repo.DB) *User {
	return &User{
		db: db,
	}
}

type User struct {
	db repo.DB
}

func (u *User) GetHomePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetHomePageResp, error) {
	permissions, err := useDahuaPermissions(ctx, u.db)
	if err != nil {
		return nil, err
	}

	dbDevices, err := dahua.ListFatDevices(ctx, u.db, permissions)
	if err != nil {
		return nil, err
	}

	devices := make([]*rpc.GetHomePageResp_Device, 0, len(dbDevices))
	for _, v := range dbDevices {
		devices = append(devices, &rpc.GetHomePageResp_Device{
			Id:   v.ID,
			Name: v.Name,
		})
	}

	return &rpc.GetHomePageResp{
		Devices: devices,
	}, nil
}

func (u *User) GetProfilePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetProfilePageResp, error) {
	session := useAuthSession(ctx)

	user, err := u.db.AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	dbSessions, err := u.db.AuthListUserSessionsForUserAndNotExpired(ctx, repo.AuthListUserSessionsForUserAndNotExpiredParams{
		UserID: session.UserID,
		Now:    types.NewTime(time.Now()),
	})
	if err != nil {
		return nil, err
	}

	activeCutoff := time.Now().Add(-24 * time.Hour)
	sessions := make([]*rpc.GetProfilePageResp_Session, 0, len(dbSessions))
	for _, v := range dbSessions {
		sessions = append(sessions, &rpc.GetProfilePageResp_Session{
			Id:             v.ID,
			UserAgent:      v.UserAgent,
			Ip:             v.Ip,
			LastIp:         v.LastIp,
			LastUsedAtTime: timestamppb.New(v.LastUsedAt.Time),
			CreatedAtTime:  timestamppb.New(v.CreatedAt.Time),
			Active:         v.LastUsedAt.After(activeCutoff),
			Current:        v.Session == session.Session,
		})
	}

	dbGroups, err := u.db.AuthListGroupsForUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	groups := make([]*rpc.GetProfilePageResp_Group, 0, len(dbGroups))
	for _, v := range dbGroups {
		groups = append(groups, &rpc.GetProfilePageResp_Group{
			Id:           v.ID,
			Name:         v.Name,
			Description:  v.Description,
			JoinedAtTime: timestamppb.New(v.JoinedAt.Time),
		})
	}

	return &rpc.GetProfilePageResp{
		Username:      user.Username,
		Email:         user.Email,
		Admin:         session.Admin,
		CreatedAtTime: timestamppb.New(user.CreatedAt.Time),
		UpdatedAtTime: timestamppb.New(user.UpdatedAt.Time),
		Sessions:      sessions,
		Groups:        groups,
	}, nil
}

func (u *User) UpdateMyPassword(ctx context.Context, req *rpc.UpdateMyPasswordReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	dbUser, err := u.db.AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if err := auth.CheckUserPassword(dbUser.Password, req.OldPassword); err != nil {
		return nil, err
	}

	if err := auth.UpdateUserPassword(ctx, u.db, dbUser, auth.UpdateUserPasswordParams{
		NewPassword:    req.NewPassword,
		CurrentSession: session.Session,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) UpdateMyUsername(ctx context.Context, req *rpc.UpdateMyUsernameReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	dbUser, err := u.db.AuthGetUser(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	if err := auth.UpdateUserUsername(ctx, u.db, dbUser, req.NewUsername); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeAllMySessions(ctx context.Context, rCreateUpdateGroupeq *emptypb.Empty) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	err := auth.DeleteOtherUserSessions(ctx, u.db, session.UserID, session.Session)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeMySession(ctx context.Context, req *rpc.RevokeMySessionReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	if err := auth.DeleteUserSession(ctx, u.db, session.UserID, req.SessionId); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
