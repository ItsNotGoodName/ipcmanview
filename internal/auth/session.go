package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

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

// CreateUserSession creates user for admin or when signing up.
func CreateUserSession(ctx context.Context, db repo.DB, arg CreateUserSessionParams) (string, error) {
	session, err := generateSession()
	if err != nil {
		return "", err
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
		ExpiredAt:  types.NewTime(now.Add(arg.Duration)),
	}

	if arg.PreviousSession != "" {
		err := createUserSessionAndDeletePrevious(ctx, db, dbArg, arg.PreviousSession)
		if err != nil {
			return "", nil
		}
	} else {
		err := db.AuthCreateUserSession(ctx, dbArg)
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

// DeleteUserSessionBySession removes session by value and is called when signing out.
func DeleteUserSessionBySession(ctx context.Context, db repo.DB, session string) error {
	return db.AuthDeleteUserSessionBySession(ctx, session)
}

// DeleteUserSession revokes session for user.
func DeleteUserSession(ctx context.Context, db repo.DB, userID int64, sessionID int64) error {
	return db.AuthDeleteUserSessionForUser(ctx, repo.AuthDeleteUserSessionForUserParams{
		UserID: userID,
		ID:     sessionID,
	})
}

// DeleteOtherUserSessions revokes all other sessions for user.
func DeleteOtherUserSessions(ctx context.Context, db repo.DB, userID int64, currentSessionID int64) error {
	return db.AuthDeleteUserSessionForUserAndNotSession(ctx, repo.AuthDeleteUserSessionForUserAndNotSessionParams{
		UserID: userID,
		ID:     currentSessionID,
	})
}

var touchSessionLock = core.NewLockStore[int64]()
var touchSessionThrottle = time.Minute

type TouchUserSessionParams struct {
	CurrentSessionID int64
	LastUsedAt       time.Time
	LastIP           string
	IP               string
}

// TouchUserSession is called whenever a session is used.
func TouchUserSession(ctx context.Context, db repo.DB, arg TouchUserSessionParams) error {
	now := time.Now()
	if arg.LastIP == arg.IP || arg.LastUsedAt.After(now.Add(-touchSessionThrottle)) {
		return nil
	}

	unlock, err := touchSessionLock.TryLock(arg.CurrentSessionID)
	if err != nil {
		return nil
	}
	defer unlock()

	err = db.AuthUpdateUserSession(ctx, repo.AuthUpdateUserSessionParams{
		LastIp:     arg.IP,
		LastUsedAt: types.NewTime(now),
		ID:         arg.CurrentSessionID,
	})
	if err != nil {
		return err
	}

	return nil
}

func UserSessionExpired(expiredAt time.Time) bool {
	return expiredAt.Before(time.Now())
}

func WithSessionAndActor(ctx context.Context, session Session) context.Context {
	ctx = core.WithUserActor(ctx, session.UserID, session.Admin)
	return context.WithValue(ctx, sessionCtxKey, session)
}

func UseSession(ctx context.Context) (Session, bool) {
	user, ok := ctx.Value(sessionCtxKey).(Session)
	return user, ok
}
