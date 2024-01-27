package rpcserver

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
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

	dbDevices, err := u.db.ListDahuaDevicesForUser(ctx, repo.ListDahuaDevicesForUserParams{
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

func (u *User) GetProfilePage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetProfilePageResp, error) {
	authSession := useAuthSession(ctx)

	user, err := u.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}

	dbSessions, err := u.db.ListUserSessionsForUserAndNotExpired(ctx, repo.ListUserSessionsForUserAndNotExpiredParams{
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

	dbGroups, err := u.db.ListGroupsForUser(ctx, authSession.UserID)
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

	dbUser, err := u.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}
	user := dbUser.Convert()

	if err := auth.CheckUserPassword(dbUser.Password, req.OldPassword); err != nil {
		return nil, NewError(err, "Old password is invalid.").Field("oldPassword")
	}

	if _, err := auth.UpdateUser(ctx, u.db, user, req.NewPassword); err != nil {
		msg := "Failed to update password."

		if errs, ok := asValidationErrors(err); ok {
			return nil, NewError(err, msg).Validation(errs, [][2]string{
				{"newPassword", "Password"},
			})
		}

		return nil, check(err)
	}

	if err := u.db.DeleteUserSessionForUserAndNotSession(ctx, repo.DeleteUserSessionForUserAndNotSessionParams{
		UserID:  authSession.UserID,
		Session: authSession.Session,
	}); err != nil {
		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}

func (u *User) UpdateMyUsername(ctx context.Context, req *rpc.UpdateMyUsernameReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	dbUser, err := u.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}
	user := dbUser.Convert()

	user.Username = req.NewUsername

	if _, err := auth.UpdateUser(ctx, u.db, user, ""); err != nil {
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

func (u *User) RevokeAllMySessions(ctx context.Context, req *rpc.RevokeAllMySessionsReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	if err := u.db.DeleteUserSessionForUserAndNotSession(ctx, repo.DeleteUserSessionForUserAndNotSessionParams{
		UserID:  authSession.UserID,
		Session: authSession.Session,
	}); err != nil {
		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}

func (u *User) RevokeMySession(ctx context.Context, req *rpc.RevokeMySessionReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	if err := u.db.DeleteUserSessionForUser(ctx, repo.DeleteUserSessionForUserParams{
		ID:     req.SessionId,
		UserID: authSession.UserID,
	}); err != nil {
		return nil, check(err)
	}

	return &emptypb.Empty{}, nil
}

// func (u *User) ListMyGroups(ctx context.Context, req *rpc.ListMyGroupsReq) (*rpc.ListMyGroupsResp, error) {
// 	authSession := useAuthSession(ctx)
//
//
// 	return &rpc.ListMyGroupsResp{
// 		Groups: groups,
// 	}, nil
// }
