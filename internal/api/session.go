package api

import (
	"net/http"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	echo "github.com/labstack/echo/v4"
)

type SesionResp struct {
	Valid    bool   `json:"valid"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

func (s *Server) Session(c echo.Context) error {
	ctx := c.Request().Context()

	authSession, ok := auth.UseSession(ctx)
	if !ok {
		return c.JSON(http.StatusUnauthorized, SesionResp{})
	}

	return c.JSON(http.StatusOK, SesionResp{
		Valid:    true,
		Username: authSession.Username,
		Admin:    authSession.Admin,
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
	user, err := s.db.GetUserByUsernameOrEmail(ctx, req.UsernameOrEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect credentials.").WithInternal(err)
	}

	// Check password
	if err := auth.CheckUserPassword(user.Password, req.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect credentials.").WithInternal(err)
	}

	// New session
	sessionDuration := auth.DefaultSessionDuration
	if req.RememberMe {
		sessionDuration = auth.RememberMeSessionDuration
	}
	session, err := auth.NewSession(ctx, s.db, c.Request().UserAgent(), c.RealIP(), user.ID, sessionDuration)
	if err != nil {
		return err
	}

	// Save session and delete previous session if it exists
	if cookie, err := c.Cookie(auth.CookieKey); err == nil {
		if err := s.db.CreateUserSessionAndDeletePrevious(ctx, session, cookie.Value); err != nil {
			return err
		}
	} else {
		if err := s.db.CreateUserSession(ctx, session); err != nil {
			return err
		}
	}

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    session.Session,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) SessionDELETE(c echo.Context) error {
	ctx := c.Request().Context()

	// Delete session
	cookie, err := c.Cookie("session")
	if err != nil {
		return c.JSON(http.StatusOK, nil)
	}
	if err := s.db.DeleteUserSessionBySession(ctx, cookie.Value); err != nil {
		return err
	}

	// Delete cookie
	c.SetCookie(&http.Cookie{
		Name:     auth.CookieKey,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, nil)
}
