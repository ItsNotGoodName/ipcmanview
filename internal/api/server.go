package api

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	echo "github.com/labstack/echo/v4"
	"github.com/spf13/afero"
)

func NewServer(
	pub pubsub.Pub,
	db repo.DB,
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
	db          repo.DB
	dahuaStore  *dahua.Store
	dahuaFileFS afero.Fs
}

func (s *Server) Register(e *echo.Echo, m ...echo.MiddlewareFunc) {
	g := e.Group("", m...)

	// Sessson
	g.GET("/v1/session", s.Session)
	g.POST("/v1/session", s.SessionPOST)
	g.DELETE("/v1/session", s.SessionDELETE)

	// Dahua
	g.GET("/v1/dahua/devices", s.DahuaDevices)
	g.GET("/v1/dahua/devices/:id/audio", s.DahuaDevicesIDAudio)
	g.GET("/v1/dahua/devices/:id/coaxial/caps", s.DahuaDevicesIDCoaxialCaps)
	g.GET("/v1/dahua/devices/:id/coaxial/status", s.DahuaDevicesIDCoaxialStatus)
	g.GET("/v1/dahua/devices/:id/detail", s.DahuaDevicesIDDetail)
	g.GET("/v1/dahua/devices/:id/error", s.DahuaDevicesIDError)
	g.GET("/v1/dahua/devices/:id/events", s.DahuaDevicesIDEvents)
	g.GET("/v1/dahua/devices/:id/files", s.DahuaDevicesIDFiles)
	g.GET("/v1/dahua/devices/:id/licenses", s.DahuaDevicesIDLicenses)
	g.GET("/v1/dahua/devices/:id/snapshot", s.DahuaDevicesIDSnapshot)
	g.GET("/v1/dahua/devices/:id/software", s.DahuaDevicesIDSoftware)
	g.GET("/v1/dahua/devices/:id/storage", s.DahuaDevicesIDStorage)
	g.GET("/v1/dahua/devices/:id/users", s.DahuaDevicesIDUsers)
	g.GET("/v1/dahua/events", s.DahuaEvents)
	g.GET(dahua.AferoEchoRoute, s.DahuaAfero(dahua.AferoEchoRoutePrefix))
	g.GET(dahua.FileEchoRoute, s.DahuaDevicesIDFilesPath)
	g.POST("/v1/dahua/devices/:id/ptz/preset", s.DahuaDevicesIDPTZPresetPOST)
	g.POST("/v1/dahua/devices/:id/rpc", s.DahuaDevicesIDRPCPOST)
}
