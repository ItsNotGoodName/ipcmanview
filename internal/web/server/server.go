package webserver

import (
	"bytes"
	"context"
	"net/http"
	"slices"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/htmx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterMiddleware(e *echo.Echo) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: web.AssetFS(),
	}))
}

func RegisterRoutes(e *echo.Echo, w Server) {
	e.DELETE("/dahua/cameras/:id", w.DahuaCamerasIDDelete)

	e.GET("/", w.Index)
	e.GET("/dahua", w.Dahua)
	e.GET("/dahua/cameras", w.DahuaCameras)
	e.GET("/dahua/cameras/:id/update", w.DahuaCamerasUpdate)
	e.GET("/dahua/cameras/create", w.DahuaCamerasCreate)
	e.GET("/dahua/events", w.DahuaEvent)
	e.GET("/dahua/events/stream", w.DahuaEventStream)
	e.GET("/dahua/files", w.DahuaFiles)
	e.GET("/dahua/snapshots", w.DahuaSnapshots)

	e.POST("/dahua/cameras/create", w.DahuaCamerasCreatePOST)
	e.POST("/dahua/cameras/:id/update", w.DahuaCamerasUpdatePOST)
}

type Server struct {
	db         sqlc.DB
	pubSub     api.PubSub
	dahuaStore *dahua.Store
	dahuaBus   *dahua.Bus
}

func New(db sqlc.DB, pubSub api.PubSub, dahuaStore *dahua.Store, dahuaBus *dahua.Bus) Server {
	return Server{
		db:         db,
		pubSub:     pubSub,
		dahuaStore: dahuaStore,
		dahuaBus:   dahuaBus,
	}
}

func (s Server) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (s Server) Dahua(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua", nil)
}

func (s Server) DahuaEvent(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-events", nil)
}

func (s Server) DahuaFiles(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-files", nil)
}

func (s Server) DahuaEventStream(c echo.Context) error {
	// FIXME: this handler prevents Echo from gracefully shutting down because the context is not canceled when Echo.Shutdown is called.

	w := c.Response()

	w.Header().Set(echo.HeaderContentType, "text/event-stream")
	w.Header().Set(echo.HeaderCacheControl, "no-cache")
	w.Header().Set(echo.HeaderConnection, "keep-alive")

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	eventsC, err := s.pubSub.SubscribeDahuaEvents(ctx, []int64{})
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	// Send previous 10 events
	events, err := s.db.ListDahuaEvent(ctx, sqlc.ListDahuaEventParams{
		Limit: 10,
	})
	if err != nil {
		return err
	}
	slices.Reverse(events)
	for _, event := range events {
		if err := c.Echo().Renderer.Render(buf, "dahua-events", TemplateBlock{
			"event-row",
			Data{
				"Event": event,
			},
		}, c); err != nil {
			return err
		}
		w.Write(formatSSE("message", buf.String()))
		buf.Reset()
		w.Flush()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-eventsC:
			if err := c.Echo().Renderer.Render(buf, "dahua-events", TemplateBlock{
				"event-row",
				Data{
					"Event": event.Event,
				},
			}, c); err != nil {
				return err
			}
			w.Write(formatSSE("message", buf.String()))
			buf.Reset()
			w.Flush()
		}
	}
}

func (s Server) DahuaCamerasIDDelete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := s.db.DeleteDahuaCamera(c.Request().Context(), id); err != nil {
		return err
	}
	s.dahuaBus.CameraDeleted(id)

	return c.NoContent(http.StatusOK)
}

func (s Server) DahuaCameras(c echo.Context) error {
	if htmx.GetRequest(c.Request()) && !htmx.GetBoosted(c.Request()) {
		apiData, err := useDahuaAPIData(c.Request().Context(), s.db, s.dahuaStore)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-cameras", TemplateBlock{
			"htmx-api-data",
			apiData,
		})
	}

	cameras, err := s.db.ListDahuaCamera(c.Request().Context())
	if err != nil {
		return err
	}

	fileCursors, err := s.db.ListDahuaFileCursor(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras", Data{
		"Cameras":     cameras,
		"FileCursors": fileCursors,
	})
}

func (s Server) DahuaCamerasCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-cameras-create", nil)
}

func (s Server) DahuaCamerasCreatePOST(c echo.Context) error {
	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
	}
	if err := parseForm(c, &form); err != nil {
		return err
	}
	if form.Name == "" {
		form.Name = form.Address
	}
	if form.Username == "" {
		form.Username = "admin"
	}
	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	dto, err := dahua.NewDahuaCamera(0, models.DTODahuaCamera{
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location,
	})

	ctx := c.Request().Context()

	id, err := s.db.CreateDahuaCamera(ctx, sqlc.CreateDahuaCameraParams{
		Name:      form.Name,
		Username:  dto.Username,
		Password:  dto.Password,
		Address:   dto.Address,
		Location:  dto.Location,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.CreatedAt,
	}, webdahua.DefaultFileCursor())
	if err != nil {
		return err
	}
	dbCamera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		return err
	}

	s.dahuaBus.CameraCreated(webdahua.ConvertGetDahuaCameraRow(dbCamera))

	return c.Redirect(http.StatusSeeOther, "/dahua/cameras")
}

func (s Server) DahuaCamerasUpdate(c echo.Context) error {
	camera, err := useDahuaCamera(c, s.db)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras-update", Data{
		"Camera": camera,
	})
}

func (s Server) DahuaCamerasUpdatePOST(c echo.Context) error {
	camera, err := useDahuaCamera(c, s.db)
	if err != nil {
		return err
	}

	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
	}
	if err := parseForm(c, &form); err != nil {
		return err
	}
	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if form.Password == "" {
		form.Password = camera.Password
	}

	dto, err := dahua.NewDahuaCamera(0, models.DTODahuaCamera{
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location,
	})

	ctx := c.Request().Context()

	_, err = s.db.UpdateDahuaCamera(ctx, sqlc.UpdateDahuaCameraParams{
		ID:       camera.ID,
		Name:     form.Name,
		Username: dto.Username,
		Password: dto.Password,
		Address:  dto.Address,
		Location: dto.Location,
	})
	if err != nil {
		return err
	}

	dbCamera, err := s.db.GetDahuaCamera(ctx, camera.ID)
	if err != nil {
		return err
	}
	s.dahuaBus.CameraUpdated(webdahua.ConvertGetDahuaCameraRow(dbCamera))

	return c.Redirect(http.StatusSeeOther, "/dahua/cameras")
}

func (s Server) DahuaSnapshots(c echo.Context) error {
	cameras, err := s.db.ListDahuaCamera(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-snapshots", Data{
		"Cameras": cameras,
	})
}
