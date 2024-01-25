package rpcserver

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewUser(db repo.DB) *User {
	return &User{
		db: db,
	}
}

type User struct {
	db repo.DB
}

func (u *User) UpdateMyPassword(ctx context.Context, req *rpc.UpdateMyPasswordReq) (*emptypb.Empty, error) {
	authSession := useAuthSession(ctx)

	dbUser, err := u.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, check(err)
	}
	user := dbUser.Convert()

	if err := auth.CheckUserPassword(dbUser.Password, req.OldPassword); err != nil {
		return nil, NewError(err, "Failed to update password.").Field("oldPassword", fmt.Errorf("Old password is invalid."))
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
