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

			// User or public actor
			var newCtx context.Context
			if session, ok := auth.UseSession(ctx); ok {
				newCtx = core.WithUserActor(ctx, session.UserID, session.Admin)
			} else {
				newCtx = core.WithPublicActor(ctx)
			}

			c.SetRequest(r.WithContext(newCtx))
			return next(c)
		}
	}
}
