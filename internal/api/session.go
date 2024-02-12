package api

import (
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	echo "github.com/labstack/echo/v4"
)

const cookieKey = "session"

// SessionMiddleware sets the session context.
func SessionMiddleware(db sqlite.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()

			cookie, err := c.Cookie(cookieKey)
			if err != nil {
				return next(c)
			}

			// Get session
			session, err := auth.GetUserSessionForContext(ctx, db, cookie.Value)
			if err != nil {
				if repo.IsNotFound(err) {
					return next(c)
				}
				return err
			}

			// Touch session
			if err := auth.TouchUserSession(ctx, db, auth.TouchUserSessionParams{
				CurrentSessionID: session.ID,
				LastUsedAt:       session.LastUsedAt.Time,
				LastIP:           session.LastIp,
				IP:               c.RealIP(),
			}); err != nil {
				if repo.IsNotFound(err) {
					return next(c)
				}
				return err
			}

			// Set session context
			c.SetRequest(r.WithContext(auth.WithSession(ctx, auth.Session{
				SessionID: session.ID,
				UserID:    session.UserID,
				Username:  session.Username.String,
				Admin:     session.Admin,
				Disabled:  session.UsersDisabledAt.Valid,
			})))
			return next(c)
		}
	}
}

type SesionResp struct {
	Admin    bool   `json:"admin"`
	Disabled bool   `json:"disabled"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}

func (s *Server) Session(c echo.Context) error {
	ctx := c.Request().Context()

	session, ok := auth.UseSession(ctx)
	if !ok {
		return c.JSON(http.StatusUnauthorized, SesionResp{})
	}

	return c.JSON(http.StatusOK, SesionResp{
		Admin:    session.Admin,
		Disabled: session.Disabled,
		UserID:   session.UserID,
		Username: session.Username,
		Valid:    true,
	})
}

func (s *Server) SessionPOST(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse request
	var req struct {
		UsernameOrEmail string
		Password        string
		RememberMe      bool
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// Get user
	user, err := auth.GetUserByUsernameOrEmail(ctx, s.db, req.UsernameOrEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect credentials.").WithInternal(err)
	}

	// Check password
	if err := auth.CheckUserPassword(user.Password, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect credentials.").WithInternal(err)
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
		RememberMe:      req.RememberMe,
		PreviousSession: previousSession,
	})
	if err != nil {
		return err
	}

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     cookieKey,
		Value:    session,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) SessionDELETE(c echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(cookieKey)
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
