package api

import (
	"net/http"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	echo "github.com/labstack/echo/v4"
)

const cookieKey = "session"

type SesionResp struct {
	Admin    bool   `json:"admin"`
	Disabled bool   `json:"disabled"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

func (s *Server) Session(c echo.Context) error {
	ctx := c.Request().Context()

	authSession, ok := auth.UseSession(ctx)
	if !ok {
		return c.JSON(http.StatusUnauthorized, SesionResp{})
	}

	return c.JSON(http.StatusOK, SesionResp{
		Admin:    authSession.Admin,
		Disabled: authSession.Disabled,
		UserID:   authSession.UserID,
		Username: authSession.Username,
		Valid:    true,
	})
}

func (s *Server) SessionPOST(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse and normalize request
	var req struct {
		UsernameOrEmail string
		Password        string
		RememberMe      bool
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	req.UsernameOrEmail = strings.ToLower(strings.TrimSpace(req.UsernameOrEmail))

	// Get user
	user, err := s.db.AuthGetUserByUsernameOrEmail(ctx, req.UsernameOrEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect credentials.").WithInternal(err)
	}

	// Check password
	if err := auth.CheckUserPassword(user.Password, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect credentials.").WithInternal(err)
	}

	sessionDuration := auth.DefaultSessionDuration
	if req.RememberMe {
		sessionDuration = auth.RememberMeSessionDuration
	}

	previousSession := ""
	if cookie, err := c.Cookie(cookieKey); err == nil {
		previousSession = cookie.Value
	}

	// Save session and delete previous session if it exists
	session, err := auth.CreateUserSession(ctx, s.db, auth.CreateUserSessionParams{
		UserAgent:       c.Request().UserAgent(),
		IP:              c.RealIP(),
		UserID:          user.ID,
		Duration:        sessionDuration,
		PreviousSession: previousSession,
	})
	if err != nil {
		return err
	}

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    session,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) SessionDELETE(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie("session")
	if err != nil {
		return c.JSON(http.StatusOK, nil)
	}

	// Delete session
	if err := auth.DeleteUserSessionBySession(ctx, s.db, cookie.Value); err != nil {
		return err
	}

	// Delete cookie
	c.SetCookie(&http.Cookie{
		Name:     cookieKey,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, nil)
}

func SessionMiddleware(db repo.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cookie, err := c.Cookie(cookieKey)
			if err != nil {
				return next(c)
			}

			session, err := db.AuthGetUserSession(ctx, cookie.Value)
			if err != nil {
				if repo.IsNotFound(err) {
					return next(c)
				}
				return err
			}
			if auth.UserSessionExpired(session.ExpiredAt.Time) {
				return next(c)
			}

			if err := auth.TouchUserSession(ctx, db, auth.TouchUserSessionParams{
				Session:    session.Session,
				LastUsedAt: session.LastUsedAt.Time,
				LastIP:     session.LastIp,
				IP:         c.RealIP(),
			}); err != nil {
				return err
			}

			c.SetRequest(c.Request().WithContext(auth.WithSession(c.Request().Context(), auth.Session{
				Admin:     session.Admin,
				Disabled:  session.UsersDisabledAt.Valid,
				Session:   session.Session,
				SessionID: session.ID,
				UserID:    session.UserID,
				Username:  session.Username.String,
			})))
			return next(c)
		}
	}
}
