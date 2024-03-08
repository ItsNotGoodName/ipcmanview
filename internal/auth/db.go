package auth

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func GetUserSessionForContext(ctx context.Context, session string) (repo.AuthGetUserSessionForContextRow, error) {
	return app.DB.C().AuthGetUserSessionForContext(ctx, repo.AuthGetUserSessionForContextParams{
		Session: session,
		Now:     types.NewTime(time.Now()),
	})
}

func ListUserSessions(ctx context.Context) ([]repo.UserSession, error) {
	actor := core.UseActor(ctx)
	return app.DB.C().AuthListUserSessionsForUserAndNotExpired(ctx, repo.AuthListUserSessionsForUserAndNotExpiredParams{
		UserID: actor.UserID,
		Now:    types.NewTime(time.Now()),
	})
}

func GetUser(ctx context.Context) (repo.User, error) {
	actor := core.UseActor(ctx)
	return app.DB.C().AuthGetUser(ctx, actor.UserID)
}

func GetUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (repo.User, error) {
	return app.DB.C().AuthGetUserByUsernameOrEmail(ctx, strings.ToLower(strings.TrimSpace(usernameOrEmail)))
}

func ListGroups(ctx context.Context) ([]repo.AuthListGroupsForUserRow, error) {
	actor := core.UseActor(ctx)
	return app.DB.C().AuthListGroupsForUser(ctx, actor.UserID)
}
