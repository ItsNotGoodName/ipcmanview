package api

import (
	"fmt"

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

func (s *Server) Register(e *echo.Echo) {
	g := e.Group(Route)

	g.GET("/session", s.Session)
	g.POST("/session", s.SessionPOST)
	g.DELETE("/session", s.SessionDELETE)

	g.GET("/dahua/afs", s.DahuaAfero("/v1/dahua/afs"))
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
}
