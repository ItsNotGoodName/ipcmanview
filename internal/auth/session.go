package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/labstack/echo/v4"
)

const DefaultSessionDuration = 30 * 24 * time.Hour

func CreateSesssion(ctx context.Context, db repo.DB, userID int64, duration time.Duration) (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	session := base64.URLEncoding.EncodeToString(b)

	now := time.Now()
	err = db.CreateUserSession(ctx, repo.CreateUserSessionParams{
		UserID:    userID,
		Session:   session,
		CreatedAt: types.NewTime(now),
		ExpiredAt: types.NewTime(now.Add(duration)),
	})
	if err != nil {
		return "", err
	}

	return session, nil
}

func SessionMiddleware(db repo.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return next(c)
			}
			session := cookie.Value

			ctx := c.Request().Context()

			user, err := db.SessionGetUserBySession(ctx, session)
			if err != nil {
				if repo.IsNotFound(err) {
					return next(c)
				}
				return err
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, sessionUserCtxKey, SessionUser{
				UserID:   user.ID,
				Username: user.Username.String,
				Admin:    user.Admin,
				Valid:    time.Now().Before(user.ExpiredAt.Time),
			})))

			return next(c)
		}
	}
}

type SessionUser struct {
	UserID   int64
	Username string
	Admin    bool
	Valid    bool
}

var sessionUserCtxKey contextKey = contextKey{"sessionUser"}

func GetSessionUser(ctx context.Context) (SessionUser, bool) {
	user, ok := ctx.Value(sessionUserCtxKey).(SessionUser)
	return user, ok
}
