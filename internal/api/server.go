package api

import (
	"fmt"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/apiws"
	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	echo "github.com/labstack/echo/v4"
	"github.com/spf13/afero"
)

func NewServer(
	pub pubsub.Pub,
	db sqlite.DB,
	dahuaStore *dahua.Store,
	dahuaFileFS afero.Fs,
) *Server {
	return &Server{
		pub:         pub,
		db:          db,
		dahuaStore:  dahuaStore,
		dahuaFileFS: dahuaFileFS,
	}
}

type Server struct {
	pub         pubsub.Pub
	db          sqlite.DB
	dahuaStore  *dahua.Store
	dahuaFileFS afero.Fs
}

const Route = "/v1"

func DahuaAferoFileURI(name string) string {
	return "/v1/dahua/afs/" + name
}

func DahuaDeviceFileURI(deviceID int64, filePath string) string {
	return fmt.Sprintf("/v1/dahua/devices/%d/files/%s", deviceID, filePath)
}

func (s *Server) RegisterSession(e *echo.Echo, m ...echo.MiddlewareFunc) *Server {
	g := e.Group(Route, m...)

	g.GET("/session", s.Session)
	g.POST("/session", s.SessionPOST)
	g.DELETE("/session", s.SessionDELETE)

	return s
}

func (s *Server) RegisterDahua(e *echo.Echo, m ...echo.MiddlewareFunc) *Server {
	g := e.Group(Route, m...)

	g.GET("/dahua/afs/*", s.DahuaAfero("/v1/dahua/afs"))
	g.GET("/dahua/events", s.DahuaEvents)

	g.GET("/dahua/devices", s.DahuaDevices)
	g.GET("/dahua/devices/:id/audio", s.DahuaDevicesIDAudio)
	g.GET("/dahua/devices/:id/coaxial/caps", s.DahuaDevicesIDCoaxialCaps)
	g.GET("/dahua/devices/:id/coaxial/status", s.DahuaDevicesIDCoaxialStatus)
	g.GET("/dahua/devices/:id/detail", s.DahuaDevicesIDDetail)
	g.GET("/dahua/devices/:id/error", s.DahuaDevicesIDError)
	g.GET("/dahua/devices/:id/events", s.DahuaDevicesIDEvents)
	g.GET("/dahua/devices/:id/files", s.DahuaDevicesIDFiles)
	g.GET("/dahua/devices/:id/files/*", s.DahuaDevicesIDFilesPath)
	g.GET("/dahua/devices/:id/licenses", s.DahuaDevicesIDLicenses)
	g.GET("/dahua/devices/:id/snapshot", s.DahuaDevicesIDSnapshot)
	g.GET("/dahua/devices/:id/software", s.DahuaDevicesIDSoftware)
	g.GET("/dahua/devices/:id/storage", s.DahuaDevicesIDStorage)
	g.GET("/dahua/devices/:id/users", s.DahuaDevicesIDUsers)
	g.POST("/dahua/devices/:id/ptz/preset", s.DahuaDevicesIDPTZPresetPOST)
	g.POST("/dahua/devices/:id/rpc", s.DahuaDevicesIDRPCPOST)

	return s
}

func (s *Server) RegisterWS(e *echo.Echo, m ...echo.MiddlewareFunc) {
	g := e.Group(Route, m...)

	g.GET("/ws", func(c echo.Context) error {
		w := c.Response()
		r := c.Request()
		ctx := r.Context()

		conn, err := apiws.Upgrade(w, r)
		if err != nil {
			return err
		}

		WS(ctx, conn, s.pub)

		return nil
	})
}

// ---------- Middleware

// ActorMiddleware sets the actor context from session context or token.
func ActorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			ctx := r.Context()

			if token := c.QueryParam("token"); token == core.RuntimeToken {
				// System
			} else if session, ok := auth.UseSession(ctx); ok {
				// User
				c.SetRequest(r.WithContext(core.WithUserActor(ctx, session.UserID, session.Admin)))
			} else {
				// Public
				c.SetRequest(r.WithContext(core.WithPublicActor(ctx)))
			}

			return next(c)
		}
	}
}

// RequireAuthMiddleware allows only if actor is system or session is valid.
func RequireAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Allow system
			if core.UseActor(ctx).Type == core.ActorTypeSystem {
				return next(c)
			}

			// Deny invalid session
			session, ok := auth.UseSession(ctx)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session or not signed in.")
			}
			if session.Disabled {
				return echo.NewHTTPError(http.StatusUnauthorized, "Account disabled.")
			}

			return next(c)
		}
	}
}
