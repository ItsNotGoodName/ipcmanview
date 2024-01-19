package api

import (
	"net/http"

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

func (s *Server) Register(e *echo.Echo) {
	e.GET("/v1/dahua/devices", s.DahuaDevices)
	e.GET("/v1/dahua/devices/:id/audio", s.DahuaDevicesIDAudio)
	e.GET("/v1/dahua/devices/:id/coaxial/caps", s.DahuaDevicesIDCoaxialCaps)
	e.GET("/v1/dahua/devices/:id/coaxial/status", s.DahuaDevicesIDCoaxialStatus)
	e.GET("/v1/dahua/devices/:id/detail", s.DahuaDevicesIDDetail)
	e.GET("/v1/dahua/devices/:id/error", s.DahuaDevicesIDError)
	e.GET("/v1/dahua/devices/:id/events", s.DahuaDevicesIDEvents)
	e.GET("/v1/dahua/devices/:id/files", s.DahuaDevicesIDFiles)
	e.GET("/v1/dahua/devices/:id/licenses", s.DahuaDevicesIDLicenses)
	e.GET("/v1/dahua/devices/:id/snapshot", s.DahuaDevicesIDSnapshot)
	e.GET("/v1/dahua/devices/:id/software", s.DahuaDevicesIDSoftware)
	e.GET("/v1/dahua/devices/:id/storage", s.DahuaDevicesIDStorage)
	e.GET("/v1/dahua/devices/:id/users", s.DahuaDevicesIDUsers)
	e.GET("/v1/dahua/events", s.DahuaEvents)
	e.GET(dahua.AferoEchoRoute, echo.WrapHandler(http.StripPrefix(dahua.AferoEchoRoutePrefix, http.FileServer(afero.NewHttpFs(s.dahuaFileFS)))))
	e.GET(dahua.FileEchoRoute, s.DahuaDevicesIDFilesPath)

	e.POST("/v1/dahua/devices/:id/ptz/preset", s.DahuaDevicesIDPTZPresetPOST)
	e.POST("/v1/dahua/devices/:id/rpc", s.DahuaDevicesIDRPCPOST)
}
