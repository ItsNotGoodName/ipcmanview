package rpcserver

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	"github.com/go-playground/validator/v10"
)

func NewUser(db repo.DB) *User {
	return &User{
		db: db,
	}
}

type User struct {
	db repo.DB
}

func (u *User) UpdatePassword(ctx context.Context, req *rpc.UserUpdatePasswordReq) (*rpc.UserUpdatePasswordResp, error) {
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

	return &rpc.UserUpdatePasswordResp{}, nil
}

func (u *User) UpdateUsername(ctx context.Context, req *rpc.UserUpdateUsernameReq) (*rpc.UserUpdateUsernameResp, error) {
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

	return &rpc.UserUpdateUsernameResp{}, nil
}

func (u *User) RevokeAllSessions(ctx context.Context, req *rpc.UserRevokeAllSessionsReq) (*rpc.UserRevokeAllSessionsResp, error) {
	authSession := useAuthSession(ctx)

	if err := u.db.DeleteUserSessionForUserAndNotSession(ctx, repo.DeleteUserSessionForUserAndNotSessionParams{
		UserID:  authSession.UserID,
		Session: authSession.Session,
	}); err != nil {
		return nil, NewError(err).Internal()
	}

	return &rpc.UserRevokeAllSessionsResp{}, nil
}

func (u *User) RevokeSession(ctx context.Context, req *rpc.UserRevokeSessionReq) (*rpc.UserRevokeSessionResp, error) {
	authSession := useAuthSession(ctx)

	if err := u.db.DeleteUserSessionForUser(ctx, repo.DeleteUserSessionForUserParams{
		ID:     req.SessionId,
		UserID: authSession.UserID,
	}); err != nil {
		return nil, NewError(err).Internal()
	}

	return &rpc.UserRevokeSessionResp{}, nil
}

func (u *User) ListGroup(ctx context.Context, req *rpc.UserListGroupReq) (*rpc.UserListGroupResp, error) {
	authSession := useAuthSession(ctx)

	dbGroups, err := u.db.ListGroupForUser(ctx, authSession.UserID)
	if err != nil {
		return nil, NewError(err).Internal()
	}

	groups := make([]*rpc.Group, 0, len(dbGroups))
	for _, v := range dbGroups {
		groups = append(groups, &rpc.Group{
			Id:          v.ID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return &rpc.UserListGroupResp{
		Groups: groups,
	}, nil
}
