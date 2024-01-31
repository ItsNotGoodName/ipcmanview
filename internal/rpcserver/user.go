package rpcserver

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
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
	authSession := useAuthSession(ctx)

	rows, err := u.db.DahuaListDevicePermissionLevels(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}
	dbDevices, err := u.db.DahuaListDevices(ctx, repo.FatDahuaDeviceParams{IDs: dahua.ListIDsByLevel(rows, models.DahuaPermissionLevelUser)})
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetHomePageResp{
		DeviceCount: int64(len(dbDevices)),
	}, nil
}

func (u *User) GetProfilePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetProfilePageResp, error) {
	authSession := useAuthSession(ctx)

	user, err := u.db.AuthGetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}

	dbSessions, err := u.db.AuthListUserSessionsForUserAndNotExpired(ctx, repo.AuthListUserSessionsForUserAndNotExpiredParams{
		UserID: authSession.UserID,
		Now:    types.NewTime(time.Now()),
	})
	if err != nil {
		return nil, check(err)
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
			Current:        v.Session == authSession.Session,
		})
	}

	dbGroups, err := u.db.AuthListGroupsForUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
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
		Admin:         authSession.Admin,
		CreatedAtTime: timestamppb.New(user.CreatedAt.Time),
		UpdatedAtTime: timestamppb.New(user.UpdatedAt.Time),
		Sessions:      sessions,
		Groups:        groups,
	}, nil
}

func (u *User) UpdateMyPassword(ctx context.Context, req *rpc.UpdateMyPasswordReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	dbUser, err := u.db.AuthGetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}

	if err := auth.CheckUserPassword(dbUser.Password, req.OldPassword); err != nil {
		return nil, NewError(err, "Old password is invalid.").Field("oldPassword")
	}

	if err := auth.UpdateUserPassword(ctx, u.db, dbUser, auth.UpdateUserPasswordParams{
		NewPassword:    req.NewPassword,
		CurrentSession: authSession.Session,
	}); err != nil {
		msg := "Failed to update password."

		if errs, ok := asValidationErrors(err); ok {
			return nil, NewError(err, msg).Validation(errs, [][2]string{
				{"newPassword", "Password"},
			})
		}

		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}

func (u *User) UpdateMyUsername(ctx context.Context, req *rpc.UpdateMyUsernameReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	dbUser, err := u.db.AuthGetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}

	if err := auth.UpdateUserUsername(ctx, u.db, dbUser, req.NewUsername); err != nil {
		msg := "Failed to update username."

		if errs, ok := asValidationErrors(err); ok {
			return nil, NewError(err, msg).Validation(errs, [][2]string{
				{"newUsername", "Username"},
			})
		}

		if constraintErr, ok := asConstraintError(err); ok {
			return nil, NewError(err, msg).Constraint(constraintErr, [][3]string{
				{"newUsername", "users.username", "Name already taken."},
			})
		}

		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeAllMySessions(ctx context.Context, rCreateUpdateGroupeq *emptypb.Empty) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	err := auth.DeleteOtherUserSessions(ctx, u.db, authSession.UserID, authSession.Session)
	if err != nil {
		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeMySession(ctx context.Context, req *rpc.RevokeMySessionReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	if err := auth.DeleteUserSession(ctx, u.db, authSession.UserID, req.SessionId); err != nil {
		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}
