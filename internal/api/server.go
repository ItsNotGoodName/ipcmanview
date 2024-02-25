package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	echo "github.com/labstack/echo/v4"
	"github.com/spf13/afero"
)

func NewServer(
	pub *pubsub.Pub,
	db sqlite.DB,
	bus *event.Bus,
	dahuaStore *dahua.Store,
	dahuaFileFS afero.Fs,
	mediamtxURL *url.URL,
) *Server {
	return &Server{
		pub:         pub,
		db:          db,
		bus:         bus,
		dahuaStore:  dahuaStore,
		dahuaFileFS: dahuaFileFS,
		mediamtxURL: mediamtxURL,
	}
}

type Server struct {
	pub         *pubsub.Pub
	db          sqlite.DB
	bus         *event.Bus
	dahuaStore  *dahua.Store
	dahuaFileFS afero.Fs
	mediamtxURL *url.URL
}

const Route = "/v1"

func MediamtxURI(path string) string {
	return fmt.Sprintf("%s/mediamtx/%s", Route, path)
}

func DahuaAferoFileURI(name string) string {
	return Route + "/dahua/afs/" + name
}

func DahuaDeviceFileURI(deviceID int64, filePath string) string {
	return fmt.Sprintf("%s/dahua/devices/%d/files/%s", Route, deviceID, filePath)
}

func (s *Server) RegisterSession(e *echo.Group) *Server {
	e.GET("/session", s.Session)
	e.POST("/session", s.SessionPOST)
	e.DELETE("/session", s.SessionDELETE)
	return s
}

func (s *Server) Register(e *echo.Group) *Server {
	e.GET("/ws", s.WS)

	e.Any("/mediamtx/*", s.Mediamtx(Route+"/mediamtx"))

	e.GET("/dahua/afs/*", s.DahuaAfero(Route+"/dahua/afs"))
	e.GET("/dahua/events", s.DahuaEvents)

	e.GET("/dahua/devices", s.DahuaDevices)
	e.GET("/dahua/devices/:id/audio", s.DahuaDevicesIDAudio)
	e.GET("/dahua/devices/:id/coaxial/caps", s.DahuaDevicesIDCoaxialCaps)
	e.GET("/dahua/devices/:id/coaxial/status", s.DahuaDevicesIDCoaxialStatus)
	e.GET("/dahua/devices/:id/detail", s.DahuaDevicesIDDetail)
	e.GET("/dahua/devices/:id/error", s.DahuaDevicesIDError)
	e.GET("/dahua/devices/:id/events", s.DahuaDevicesIDEvents)
	e.GET("/dahua/devices/:id/files", s.DahuaDevicesIDFiles)
	e.GET("/dahua/devices/:id/files/*", s.DahuaDevicesIDFilesPath)
	e.GET("/dahua/devices/:id/licenses", s.DahuaDevicesIDLicenses)
	e.GET("/dahua/devices/:id/ptz/preset", s.DahuaDevicesIDPTZPresetGET)
	e.GET("/dahua/devices/:id/snapshot", s.DahuaDevicesIDSnapshot)
	e.GET("/dahua/devices/:id/software", s.DahuaDevicesIDSoftware)
	e.GET("/dahua/devices/:id/storage", s.DahuaDevicesIDStorage)
	e.GET("/dahua/devices/:id/users", s.DahuaDevicesIDUsers)

	e.POST("/dahua/devices/:id/ptz/preset", s.DahuaDevicesIDPTZPresetPOST)
	e.POST("/dahua/devices/:id/rpc", s.DahuaDevicesIDRPCPOST)

	return s
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
