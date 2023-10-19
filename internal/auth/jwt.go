package auth

import (
	"context"
	"crypto/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/rpcgen"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func JWTRandomSecret() []byte {
	secret := make([]byte, 64)
	rand.Read(secret)
	return secret
}

type JWTClaim struct {
	UserID    int64
	ExpiredAt time.Time
}

func NewJWTClaim(userID int64, duration time.Duration) JWTClaim {
	return JWTClaim{
		UserID:    userID,
		ExpiredAt: time.Now().Add(duration),
	}
}

type JWTAuth struct {
	jwtAuth  *jwtauth.JWTAuth
	Verifier func(http.Handler) http.Handler
}

func NewJWTAuthenticator(secret []byte) JWTAuth {
	jwtAuth := jwtauth.New("HS256", secret, nil)
	return JWTAuth{
		jwtAuth:  jwtAuth,
		Verifier: jwtauth.Verifier(jwtAuth),
	}
}

func (j JWTAuth) Encode(claim JWTClaim) (string, error) {
	e := map[string]interface{}{
		"user_id": strconv.FormatInt(claim.UserID, 16),
		"exp":     claim.ExpiredAt.Unix(),
	}
	_, string, err := j.jwtAuth.Encode(e)
	return string, err
}

var jwtClaimCtxKey contextKey = contextKey{"claim"}

func JWTClaimFromContext(ctx context.Context) JWTClaim {
	token, _ := ctx.Value(jwtClaimCtxKey).(JWTClaim)
	return token
}

// JWTAuthenticator validates and parses JWT.
func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get token and claims
		token, claims, err := jwtauth.FromContext(ctx)
		if err != nil {
			if rpcgen.IsRPC(r) {
				rpcgen.RespondWithError(w, rpcgen.ErrInvalidToken)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}
		if token == nil || jwt.Validate(token) != nil {
			if rpcgen.IsRPC(r) {
				rpcgen.RespondWithError(w, rpcgen.ErrInvalidToken)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}

		// Parse claims
		userID, err := strconv.ParseInt(string(claims["user_id"].(string)), 10, 64)
		if err != nil {
			if rpcgen.IsRPC(r) {
				rpcgen.RespondWithError(w, rpcgen.ErrWebrpcInternalError)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, jwtClaimCtxKey, JWTClaim{
			UserID:    userID,
			ExpiredAt: token.Expiration(),
		})))
	})
}
