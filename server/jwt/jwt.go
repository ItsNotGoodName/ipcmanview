package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/ipcmango/server/service"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			service.RespondWithError(w, service.ErrInvalidToken)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			service.RespondWithError(w, service.ErrInvalidToken)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func DecodeUserID(ctx context.Context) int64 {
	_, claims, _ := jwtauth.FromContext(ctx)
	fmt.Printf("%+v", claims["user_id"])
	id, _ := strconv.ParseInt(string(claims["user_id"].(string)), 10, 64)
	return id
}

func EncodeUserID(userID int64) string {
	e := map[string]interface{}{"user_id": strconv.FormatInt(userID, 16)}
	_, string, _ := TokenAuth.Encode(e)
	return string
}
