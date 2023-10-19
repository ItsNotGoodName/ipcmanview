package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/db"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcgen"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/jackc/pgx/v5"
)

const AuthAccessTokenDuration = 5 * time.Minute
const AuthSessionTokenDuration = 30 * 24 * time.Hour

type AuthService struct {
	jwt auth.JWTAuth
	db  qes.Querier
}

var _ rpcgen.AuthService = (*AuthService)(nil)

func NewAuthService(db qes.Querier, jwt auth.JWTAuth) *AuthService {
	return &AuthService{
		db:  db,
		jwt: jwt,
	}
}

// Login implements rpcgen.AuthService.
func (a *AuthService) Login(ctx context.Context, usernameOrEmail string, password string, clientId string, ipAddress string) (*rpcgen.AuthLoginResponse, error) {
	// Get user
	user, err := db.User.GetByUsernameOrEmail(ctx, a.db, usernameOrEmail)
	if err != nil {
		return nil, handleErr(rpcgen.ErrWebrpcBadRequest, err)
	}

	// Verify password
	if err := user.CheckPassword(password); err != nil {
		return nil, handleErr(rpcgen.ErrWebrpcBadRequest, err)
	}

	// Create refresh token
	newSession, err := auth.NewSession(clientId, ipAddress, user.ID, AuthSessionTokenDuration)
	if err != nil {
		return nil, handleErr(rpcgen.ErrWebrpcBadRequest, err)
	}

	// Save refresh token
	session, err := auth.DB.SessionCreate(ctx, a.db, newSession)
	if err != nil {
		return nil, handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	// Create access token
	accessToken, err := a.jwt.Encode(auth.NewJWTClaim(user.ID, AuthAccessTokenDuration))
	if err != nil {
		return nil, handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	return &rpcgen.AuthLoginResponse{
		Session:   session.Token,
		Token:     accessToken,
		ExpiredAt: session.ExpiredAt,
	}, nil
}

// Refresh implements rpcgen.AuthService.
func (a *AuthService) Refresh(ctx context.Context, sessionToken string) (string, error) {
	// Get session
	session, err := auth.DB.SessionGet(ctx, a.db, sessionToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", handleErr(rpcgen.ErrInvalidSession, err)
		}
		return "", handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	// Validate session
	if !session.Valid() {
		if err := auth.DB.SessionRemove(ctx, a.db, session.Token); err != nil {
			return "", handleErr(rpcgen.ErrWebrpcInternalError, err)
		}

		return "", rpcgen.ErrInvalidSession
	}

	// Create access token
	accessToken, err := a.jwt.Encode(auth.NewJWTClaim(session.UserID, AuthAccessTokenDuration))
	if err != nil {
		return "", handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	return accessToken, nil
}

// Register implements rpcgen.AuthService.
func (a *AuthService) Register(ctx context.Context, email string, username string, password string, passwordConfirm string) error {
	// Create user
	user, err := core.NewUser(core.UserCreate{
		Username:        username,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	})
	if err != nil {
		return handleErr(rpcgen.ErrWebrpcBadRequest, err)
	}

	// Save user
	_, err = db.User.Create(ctx, a.db, user)
	if err != nil {
		return handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	return nil
}

// Logout implements rpcgen.AuthService.
func (a *AuthService) Logout(ctx context.Context, sessionToken string) error {
	// Remove session
	if err := auth.DB.SessionRemove(ctx, a.db, sessionToken); err != nil {
		return handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	return nil
}
