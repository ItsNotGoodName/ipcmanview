package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func generateSession() (string, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

type Session struct {
	SessionID int64
	UserID    int64
	Username  string
	Admin     bool
	Disabled  bool
}

type sessionCtxKey struct{}

const defaultSessionDuration = 24 * time.Hour          // 1 Day
const rememberMeSessionDuration = 365 * 24 * time.Hour // 1 Year

type CreateUserSessionParams struct {
	UserAgent       string
	IP              string
	UserID          int64
	RememberMe      bool
	PreviousSession string
}

func CreateUserSession(ctx context.Context, arg CreateUserSessionParams) (string, error) {
	session, err := generateSession()
	if err != nil {
		return "", err
	}

	sessionDuration := defaultSessionDuration
	if arg.RememberMe {
		sessionDuration = rememberMeSessionDuration
	}
	now := time.Now()
	dbArg := repo.AuthCreateUserSessionParams{
		UserID:     arg.UserID,
		Session:    session,
		UserAgent:  arg.UserAgent,
		Ip:         arg.IP,
		LastIp:     arg.IP,
		LastUsedAt: types.NewTime(now),
		CreatedAt:  types.NewTime(now),
		ExpiredAt:  types.NewTime(now.Add(sessionDuration)),
	}

	if arg.PreviousSession != "" {
		err := createUserSessionAndDeletePrevious(ctx, dbArg, arg.PreviousSession)
		if err != nil {
			return "", err
		}
	} else {
		err := app.DB.C().AuthCreateUserSession(ctx, dbArg)
		if err != nil {
			return "", err
		}
	}

	return session, nil
}

func createUserSessionAndDeletePrevious(ctx context.Context, arg repo.AuthCreateUserSessionParams, previousSession string) error {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.C().AuthCreateUserSession(ctx, arg); err != nil {
		return err
	}

	if _, err := tx.C().AuthDeleteUserSessionBySession(ctx, previousSession); err != nil && !core.IsNotFound(err) {
		return err
	}

	return tx.Commit()
}

func DeleteUserSessionBySession(ctx context.Context, session string) error {
	userID, err := app.DB.C().AuthDeleteUserSessionBySession(ctx, session)
	if err != nil {
		if core.IsNotFound(err) {
			return nil
		}
		return err
	}

	app.Hub.UserSecurityUpdated(bus.UserSecurityUpdated{
		UserID: userID,
	})

	return nil
}

func DeleteUserSession(ctx context.Context, userID int64, sessionID int64) error {
	if _, err := core.AssertAdminOrUser(ctx, userID); err != nil {
		return err
	}

	err := app.DB.C().AuthDeleteUserSessionForUser(ctx, repo.AuthDeleteUserSessionForUserParams{
		UserID: userID,
		ID:     sessionID,
	})
	if err != nil {
		return err
	}

	app.Hub.UserSecurityUpdated(bus.UserSecurityUpdated{
		UserID: userID,
	})

	return nil
}

func DeleteOtherUserSessions(ctx context.Context, userID int64, currentSessionID int64) error {
	if _, err := core.AssertAdminOrUser(ctx, userID); err != nil {
		return err
	}

	err := app.DB.C().AuthDeleteUserSessionForUserAndNotSession(ctx, repo.AuthDeleteUserSessionForUserAndNotSessionParams{
		UserID: userID,
		ID:     currentSessionID,
	})
	if err != nil {
		return err
	}

	app.Hub.UserSecurityUpdated(bus.UserSecurityUpdated{
		UserID: userID,
	})

	return nil
}

func NewTouchSessionThrottle() TouchSessionThrottle {
	return TouchSessionThrottle{
		LockStore: core.NewLockStore[int64](),
		Duration:  time.Minute,
	}
}

type TouchSessionThrottle struct {
	*core.LockStore[int64]
	Duration time.Duration
}

type TouchUserSessionParams struct {
	CurrentSessionID int64
	LastUsedAt       time.Time
	LastIP           string
	IP               string
}

func TouchUserSession(ctx context.Context, arg TouchUserSessionParams) error {
	now := time.Now()
	if arg.LastIP == arg.IP && arg.LastUsedAt.After(now.Add(-app.TouchSessionThrottle.Duration)) {
		return nil
	}

	unlock, err := app.TouchSessionThrottle.TryLock(arg.CurrentSessionID)
	if err != nil {
		return nil
	}
	defer unlock()

	err = app.DB.C().AuthUpdateUserSession(ctx, repo.AuthUpdateUserSessionParams{
		LastIp:     arg.IP,
		LastUsedAt: types.NewTime(now),
		ID:         arg.CurrentSessionID,
	})
	if err != nil {
		return err
	}

	return nil
}

func WithSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionCtxKey{}, session)
}

func UseSession(ctx context.Context) (Session, bool) {
	user, ok := ctx.Value(sessionCtxKey{}).(Session)
	return user, ok
}
