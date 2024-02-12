package api

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	echo "github.com/labstack/echo/v4"
)

func ActorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()

			var newCtx context.Context
			if token := c.QueryParam("token"); token == core.RuntimeToken {
				// System
			} else if session, ok := auth.UseSession(ctx); ok {
				// User
				newCtx = core.WithUserActor(ctx, session.UserID, session.Admin)
			} else {
				// Public
				newCtx = core.WithPublicActor(ctx)
			}

			c.SetRequest(r.WithContext(newCtx))
			return next(c)
		}
	}
}
