package http

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/pkg/echoext"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	ads *api.DahuaServer,
) *echo.Echo {
	e := echo.New()
	echoext.WithErrorLogging(e)

	// Middleware
	e.Use(echoext.LoggerWithConfig(echoext.LoggerConfig{
		Format: []string{
			"remote_ip",
			"host",
			"method",
			"user_agent",
			"status",
			"error",
			"latency_human",
		},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: echoext.RecoverLogErrorFunc,
	}))

	// Routes
	e.GET("/v1/dahua", ads.GET)
	e.POST("/v1/dahua", ads.POST)
	e.POST("/v1/dahua/:id", ads.IDPOST)
	e.GET("/v1/dahua-events", ads.EventsGET)
	e.POST("/v1/dahua/:id/rpc", ads.IDRPCPOST)
	e.GET("/v1/dahua/:id/detail", ads.IDDetailGET)
	e.GET("/v1/dahua/:id/software", ads.IDSoftwareGET)
	e.GET("/v1/dahua/:id/licenses", ads.IDLicensesGET)
	e.GET("/v1/dahua/:id/error", ads.IDErrorGET)
	e.GET("/v1/dahua/:id/snapshot", ads.IDSnapshotGET)
	e.GET("/v1/dahua/:id/events", ads.IDEventsGET)
	e.GET("/v1/dahua/:id/files", ads.IDFilesGET)
	e.GET("/v1/dahua/:id/files/*", ads.IDFilesPathGET)
	e.GET("/v1/dahua/:id/audio", ads.IDAudioGET)
	e.GET("/v1/dahua/:id/coaxial/status", ads.IDCoaxialStatusGET)
	e.GET("/v1/dahua/:id/coaxial/caps", ads.IDCoaxialCapsGET)
	e.POST("/v1/dahua/:id/ptz/preset", ads.IDPTZPresetPOST)
	e.GET("/v1/dahua/:id/storage", ads.IDStorageGET)
	e.GET("/v1/dahua/:id/users", ads.IDUsersGET)

	return e
}
