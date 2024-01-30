package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Session struct {
	Admin     bool
	Disabled  bool
	Session   string
	SessionID int64
	UserID    int64
	Username  string
}

var sessionCtxKey contextKey = contextKey("session")

const DefaultSessionDuration = 24 * time.Hour          // 1 Day
const RememberMeSessionDuration = 365 * 24 * time.Hour // 1 Year

type CreateUserSessionParams struct {
	UserAgent       string
	IP              string
	UserID          int64
	Duration        time.Duration
	PreviousSession string
}

func CreateUserSession(ctx context.Context, db repo.DB, arg CreateUserSessionParams) (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	session := base64.URLEncoding.EncodeToString(b)

	now := time.Now()
	args := repo.AuthCreateUserSessionParams{
		UserID:     arg.UserID,
		Session:    session,
		UserAgent:  arg.UserAgent,
		Ip:         arg.IP,
		LastIp:     arg.IP,
		LastUsedAt: types.NewTime(now),
		CreatedAt:  types.NewTime(now),
		ExpiredAt:  types.NewTime(now.Add(arg.Duration)),
	}

	if arg.PreviousSession != "" {
		err := createUserSessionAndDeletePrevious(ctx, db, args, arg.PreviousSession)
		if err != nil {
			return "", nil
		}
	} else {
		err := db.AuthCreateUserSession(ctx, args)
		if err != nil {
			return "", nil
		}
	}

	return session, nil
}

func createUserSessionAndDeletePrevious(ctx context.Context, db repo.DB, arg repo.AuthCreateUserSessionParams, previousSession string) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.AuthCreateUserSession(ctx, arg); err != nil {
		return err
	}

	if err := tx.AuthDeleteUserSessionBySession(ctx, previousSession); err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteUserSession(ctx context.Context, db repo.DB, session string) error {
	return db.AuthDeleteUserSessionBySession(ctx, session)
}

const CookieKey = "session"

func SessionMiddleware(db repo.DB) echo.MiddlewareFunc {
	sessionUpdateLock := core.NewLockStore[string]()
	sessionUpdateThrottle := time.Minute
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cookie, err := c.Cookie(CookieKey)
			if err != nil {
				return next(c)
			}

			// Get valid user
			userSession, err := db.AuthGetUserBySession(ctx, cookie.Value)
			if err != nil {
				if repo.IsNotFound(err) {
					return next(c)
				}
				return err
			}
			if userSession.ExpiredAt.Time.Before(time.Now()) {
				return next(c)
			}

			// Update last used at and last ip
			realIP := c.RealIP()
			now := time.Now()
			if userSession.LastIp != realIP || userSession.LastUsedAt.Before(now.Add(-sessionUpdateThrottle)) {
				unlock, err := sessionUpdateLock.TryLock(cookie.Value)
				if err == nil {
					err := db.AuthUpdateUserSession(ctx, repo.AuthUpdateUserSessionParams{
						LastIp:     realIP,
						LastUsedAt: types.NewTime(now),
						Session:    cookie.Value,
					})
					if err != nil {
						log.Err(err).Send()
					}
					unlock()
				}
			}

			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, sessionCtxKey, Session{
				Admin:     userSession.Admin,
				Disabled:  userSession.UsersDisabledAt.Valid,
				Session:   cookie.Value,
				SessionID: userSession.ID,
				UserID:    userSession.UserID,
				Username:  userSession.Username.String,
			})))
			return next(c)
		}
	}
}

// UseSession gets user from context.
// It fails when session does not exist or is invalid.
func UseSession(ctx context.Context) (Session, bool) {
	user, ok := ctx.Value(sessionCtxKey).(Session)
	return user, ok
}
