package rpcserver

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/go-playground/validator/v10"
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

func (u *User) UpdateMyPassword(ctx context.Context, req *rpc.UpdateMyPasswordReq) (*rpc.UpdateMyPasswordResp, error) {
	authSession := useAuthSession(ctx)

	dbUser, err := u.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, NewError(err).Internal()
	}
	user := dbUser.Convert()

	if err := auth.CheckUserPassword(dbUser.Password, req.OldPassword); err != nil {
		return nil, NewError(err, "Failed to update password.").Field("oldPassword", fmt.Errorf("Old password is invalid."))
	}

	if _, err := auth.UpdateUser(ctx, u.db, user, req.NewPassword); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return nil, NewError(err, "Failed to update password.").Validation(errs, [][2]string{
				{"Password", "newPassword"},
			})
		}

		return nil, NewError(err).Internal()
	}

	if err := u.db.DeleteUserSessionForUserAndNotSession(ctx, repo.DeleteUserSessionForUserAndNotSessionParams{
		UserID:  authSession.UserID,
		Session: authSession.Session,
	}); err != nil {
		return nil, NewError(err).Internal()
	}

	return &rpc.UpdateMyPasswordResp{}, nil
}

func (u *User) UpdateMyUsername(ctx context.Context, req *rpc.UpdateMyUsernameReq) (*rpc.UpdateMyUsernameResp, error) {
	authSession := useAuthSession(ctx)

	dbUser, err := u.db.GetUser(ctx, authSession.UserID)
	if err != nil {
		return nil, NewError(err).Internal()
	}
	user := dbUser.Convert()

	user.Username = req.NewUsername

	if _, err := auth.UpdateUser(ctx, u.db, user, ""); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return nil, NewError(err, "Failed to update username.").Validation(errs, [][2]string{
				{"Username", "newUsername"},
			})
		}

		if constraintErr, ok := sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE); ok {
			return nil, NewError(err, "Failed to update username.").Constraint(constraintErr, [][3]string{
				{"users.username", "newUsername", "Name already taken."},
			})
		}

		return nil, NewError(err).Internal()
	}

	return &rpc.UpdateMyUsernameResp{}, nil
}

func (u *User) RevokeAllMySessions(ctx context.Context, req *rpc.RevokeAllMySessionsReq) (*rpc.RevokeAllMySessionsResp, error) {
	authSession := useAuthSession(ctx)

	if err := u.db.DeleteUserSessionForUserAndNotSession(ctx, repo.DeleteUserSessionForUserAndNotSessionParams{
		UserID:  authSession.UserID,
		Session: authSession.Session,
	}); err != nil {
		return nil, NewError(err).Internal()
	}

	return &rpc.RevokeAllMySessionsResp{}, nil
}

func (u *User) RevokeMySession(ctx context.Context, req *rpc.RevokeMySessionReq) (*rpc.RevokeMySessionResp, error) {
	authSession := useAuthSession(ctx)

	if err := u.db.DeleteUserSessionForUser(ctx, repo.DeleteUserSessionForUserParams{
		ID:     req.SessionId,
		UserID: authSession.UserID,
	}); err != nil {
		return nil, NewError(err).Internal()
	}

	return &rpc.RevokeMySessionResp{}, nil
}

func (u *User) ListMyGroups(ctx context.Context, req *rpc.ListMyGroupsReq) (*rpc.ListMyGroupsResp, error) {
	authSession := useAuthSession(ctx)

	dbGroups, err := u.db.ListGroupForUser(ctx, authSession.UserID)
	if err != nil {
		return nil, NewError(err).Internal()
	}

	groups := make([]*rpc.Group, 0, len(dbGroups))
	for _, v := range dbGroups {
		groups = append(groups, &rpc.Group{
			Id:           v.ID,
			Name:         v.Name,
			Description:  v.Description,
			JoinedAtTime: timestamppb.New(v.JoinedAt.Time),
		})
	}

	return &rpc.ListMyGroupsResp{
		Groups: groups,
	}, nil
}
