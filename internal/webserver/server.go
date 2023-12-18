package webserver

import (
	"bytes"
	"context"
	"net/http"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/pkg/htmx"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func RegisterMiddleware(e *echo.Echo) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: web.AssetFS(),
	}))
}

func RegisterRoutes(e *echo.Echo, w Server) {
	e.GET("/", w.Index)
	e.GET("/dahua", w.Dahua)
	e.GET("/dahua/cameras", w.DahuaCameras)
	e.GET("/dahua/cameras/:id/update", w.DahuaCamerasUpdate)
	e.GET("/dahua/cameras/create", w.DahuaCamerasCreate)
	e.GET("/dahua/cameras/file-cursors", w.DahuaCamerasFileCursors)
	e.GET("/dahua/events", w.DahuaEvents)
	e.GET("/dahua/events/:id/data", w.DahuaEventsIDData)
	e.GET("/dahua/events/live", w.DahuaEventsLive)
	e.GET("/dahua/events/rules", w.DahuaEventsRules)
	e.GET("/dahua/events/stream", w.DahuaEventStream)
	e.GET("/dahua/files", w.DahuaFiles)
	e.GET("/dahua/snapshots", w.DahuaSnapshots)

	e.POST("/dahua/cameras", w.DahuaCamerasPOST)
	e.POST("/dahua/cameras/:id/update", w.DahuaCamerasUpdatePOST)
	e.POST("/dahua/cameras/create", w.DahuaCamerasCreatePOST)
	e.POST("/dahua/cameras/file-cursors", w.DahuaCamerasFileCursorsPOST)
	e.POST("/dahua/events/rules", w.DahuaEventsRulePOST)
	e.POST("/dahua/events/rules/create", w.DahuaEventsRulesCreatePOST)
}

type Server struct {
	db         repo.DB
	pub        pubsub.Pub
	bus        *core.Bus
	dahuaStore *dahuacore.Store
}

func New(db repo.DB, pub pubsub.Pub, bus *core.Bus, dahuaStore *dahuacore.Store) Server {
	return Server{
		db:         db,
		pub:        pub,
		bus:        bus,
		dahuaStore: dahuaStore,
	}
}

func (s Server) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (s Server) Dahua(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua", nil)
}

func (s Server) DahuaEvents(c echo.Context) error {
	ctx := c.Request().Context()

	params := struct {
		CameraID  []int64
		Code      []string
		Action    []string
		Page      int64 `validate:"gt=0"`
		PerPage   int64
		Data      bool
		Start     string
		End       string
		Ascending bool
	}{
		Page:    1,
		PerPage: 10,
	}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}
	if err := api.ValidateStruct(params); err != nil {
		return err
	}

	timeRange, err := api.UseTimeRange(params.Start, params.End)
	if err != nil {
		return err
	}

	events, err := s.db.ListDahuaEvent(ctx, repo.ListDahuaEventParams{
		CameraID: params.CameraID,
		Code:     params.Code,
		Action:   params.Action,
		Page: pagination.Page{
			Page:    int(params.Page),
			PerPage: int(params.PerPage),
		},
		Start:     types.NewTime(timeRange.Start),
		End:       types.NewTime(timeRange.End),
		Ascending: params.Ascending,
	})
	if err != nil {
		return err
	}

	eventCodes, err := s.db.ListDahuaEventCodes(ctx)
	if err != nil {
		return err
	}

	eventActions, err := s.db.ListDahuaEventActions(ctx)
	if err != nil {
		return err
	}

	cameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Params":      params,
		"Cameras":     cameras,
		"Events":      events,
		"EventCodes":  eventCodes,
		"EventAction": eventActions,
	}

	if isHTMX(c) {
		htmx.SetReplaceURL(c.Response(), c.Request().URL.Path+"?"+api.EncodeQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-events", TemplateBlock{"htmx", data})
	}

	return c.Render(http.StatusOK, "dahua-events", data)
}

func (s Server) DahuaEventsIDData(c echo.Context) error {
	id, err := api.ParamID(c)
	if err != nil {
		return err
	}

	data, err := s.db.GetDahuaEventData(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-events", TemplateBlock{"dahua-events-data-json", data})
}

func (s Server) DahuaFiles(c echo.Context) error {
	ctx := c.Request().Context()

	params := struct {
		CameraID  []int64
		Type      []string
		Page      int64 `validate:"gt=0"`
		PerPage   int64
		Start     string
		End       string
		Ascending bool
	}{
		Page:    1,
		PerPage: 10,
	}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}
	if err := api.ValidateStruct(params); err != nil {
		return err
	}

	timeRange, err := api.UseTimeRange(params.Start, params.End)
	if err != nil {
		return err
	}

	files, err := s.db.ListDahuaFile(ctx, repo.ListDahuaFileParams{
		Page: pagination.Page{
			Page:    int(params.Page),
			PerPage: int(params.PerPage),
		},
		Type:      params.Type,
		CameraID:  params.CameraID,
		Start:     types.NewTime(timeRange.Start),
		End:       types.NewTime(timeRange.End),
		Ascending: params.Ascending,
	})
	if err != nil {
		return err
	}

	cameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	types, err := s.db.ListDahuaFileTypes(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Params":  params,
		"Files":   files,
		"Cameras": cameras,
		"Types":   types,
	}

	if isHTMX(c) {
		htmx.SetReplaceURL(c.Response(), c.Request().URL.Path+"?"+api.EncodeQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-files", TemplateBlock{"htmx", data})
	}

	return c.Render(http.StatusOK, "dahua-files", data)
}

func (s Server) DahuaEventsLive(c echo.Context) error {
	ctx := c.Request().Context()

	params := struct {
		CameraID []int64
		Code     []string
		Action   []string
		Data     bool
	}{}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}

	eventCodes, err := s.db.ListDahuaEventCodes(ctx)
	if err != nil {
		return err
	}

	eventActions, err := s.db.ListDahuaEventActions(ctx)
	if err != nil {
		return err
	}

	cameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Params":      params,
		"Cameras":     cameras,
		"EventCodes":  eventCodes,
		"EventAction": eventActions,
	}

	if isHTMX(c) {
		htmx.SetReplaceURL(c.Response(), c.Request().URL.Path+"?"+api.EncodeQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-events-live", TemplateBlock{"htmx", data})
	}

	return c.Render(http.StatusOK, "dahua-events-live", data)
}

func (s Server) DahuaEventStream(c echo.Context) error {
	ctx := c.Request().Context()

	params := struct {
		CameraID []int64
		Code     []string
		Action   []string
		Data     bool
	}{}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}

	sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaCameraEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	w := useEventStream(c)
	buf := new(bytes.Buffer)

	for event := range eventsC {
		evt, ok := event.(models.EventDahuaCameraEvent)
		if !ok ||
			evt.EventRule.IgnoreLive ||
			(len(params.CameraID) > 0 && !slices.Contains(params.CameraID, evt.Event.CameraID)) ||
			(len(params.Code) > 0 && !slices.Contains(params.Code, evt.Event.Code)) ||
			(len(params.Action) > 0 && !slices.Contains(params.Action, evt.Event.Action)) {
			continue
		}

		if err := c.Echo().Renderer.Render(buf, "dahua-events-live", TemplateBlock{"event-row", Data{
			"Event":  evt.Event,
			"Params": params,
		}}, c); err != nil {
			return err
		}

		err := sendEventStream(w, formatEventStream("message", buf.String()))
		if err != nil {
			return err
		}

		buf.Reset()
	}

	return sub.Error()
}

func (s Server) DahuaCamerasPOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Action  string
		Cameras []struct {
			Selected bool
			ID       int64
		}
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	if form.Action == "Delete" {
		for _, camera := range form.Cameras {
			if !camera.Selected {
				continue
			}
			if err := dahua.DeleteCamera(ctx, s.db, s.bus, camera.ID); err != nil {
				return err
			}
		}
	}

	if isHTMX(c) {
		cameras, err := s.db.ListDahuaCamera(c.Request().Context())
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-cameras", TemplateBlock{"htmx-cameras", Data{
			"Cameras": cameras,
		}})
	}

	return s.DahuaCameras(c)
}

func (s Server) DahuaCameras(c echo.Context) error {
	if isHTMX(c) {
		tables, err := useDahuaTables(c.Request().Context(), s.db, s.dahuaStore)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-cameras", TemplateBlock{"htmx", tables})
	}

	cameras, err := s.db.ListDahuaCamera(c.Request().Context())
	if err != nil {
		return err
	}

	fileCursors, err := s.db.ListDahuaFileCursor(c.Request().Context(), dahua.ScanLockStaleTime())
	if err != nil {
		return err
	}

	eventWorkers, err := s.db.ListDahuaEventWorkerState(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras", Data{
		"Cameras":      cameras,
		"FileCursors":  fileCursors,
		"EventWorkers": eventWorkers,
	})
}

func (s Server) DahuaCamerasCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-cameras-create", Data{
		"Locations": dahua.Locations,
	})
}

func (s Server) DahuaCamerasCreatePOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
	}
	if err := api.ParseForm(c, &form); err != nil {
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

	create, err := dahuacore.NewDahuaCamera(models.DahuaCamera{
		Name:     form.Name,
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location.Location,
	})

	id, err := s.db.CreateDahuaCamera(ctx, repo.CreateDahuaCameraParams{
		Name:      create.Name,
		Username:  create.Username,
		Password:  create.Password,
		Address:   create.Address,
		Location:  types.NewLocation(create.Location),
		CreatedAt: types.NewTime(create.CreatedAt),
		UpdatedAt: types.NewTime(create.UpdatedAt),
	}, dahua.NewFileCursor())
	if err != nil {
		return err
	}
	dbCamera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		return err
	}
	s.bus.EventDahuaCameraCreated(models.EventDahuaCameraCreated{
		Camera: dbCamera.Convert(),
	})

	return c.Redirect(http.StatusSeeOther, "/dahua/cameras")
}

func (s Server) DahuaCamerasFileCursorsPOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Action      string
		FileCursors []struct {
			Selected bool
			CameraID int64
		}
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	switch form.Action {
	case "Reset":
		for _, v := range form.FileCursors {
			if !v.Selected {
				continue
			}

			err := dahua.ScanLockCreate(ctx, s.db, v.CameraID)
			if err != nil {
				return err
			}

			err = dahua.ScanReset(ctx, s.db, v.CameraID)

			dahua.ScanLockDelete(s.db, v.CameraID)

			if err != nil {
				if repo.IsNotFound(err) {
					continue
				}
				return err
			}
		}
	case "Quick", "Full":
		scanType := dahua.ScanTypeQuick
		if form.Action == "Full" {
			scanType = dahua.ScanTypeFull
		}

		for _, v := range form.FileCursors {
			if !v.Selected {
				continue
			}

			camera, err := s.db.GetDahuaCamera(ctx, v.CameraID)
			if err != nil {
				if repo.IsNotFound(err) {
					continue
				}
				return err
			}
			conn := s.dahuaStore.Conn(ctx, camera.Convert().DahuaConn)

			if err := dahua.ScanLockCreate(ctx, s.db, v.CameraID); err != nil {
				return err
			}
			go func(conn dahuacore.Conn) {
				ctx := context.Background()
				cancel := dahua.ScanLockHeartbeat(ctx, s.db, conn.Camera.ID)
				defer cancel()

				err := dahua.Scan(ctx, s.db, conn.RPC, conn.Camera, scanType)
				if err != nil {
					log.Err(err).Msg("Scan error")
				}
			}(conn)
		}
	}

	return s.DahuaCamerasFileCursors(c)
}

func (s Server) DahuaCamerasFileCursors(c echo.Context) error {
	if !isHTMX(c) {
		return c.Redirect(http.StatusSeeOther, "/dahua/cameras#file-cursors")
	}

	fileCursors, err := s.db.ListDahuaFileCursor(c.Request().Context(), dahua.ScanLockStaleTime())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras", TemplateBlock{"htmx-file-cursors", Data{
		"FileCursors": fileCursors,
	}})
}

func (s Server) DahuaCamerasUpdate(c echo.Context) error {
	camera, err := useDahuaCamera(c, s.db)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-cameras-update", Data{
		"Locations": dahua.Locations,
		"Camera":    camera,
	})
}

func (s Server) DahuaCamerasUpdatePOST(c echo.Context) error {
	ctx := c.Request().Context()

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
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}
	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if form.Password == "" {
		form.Password = camera.Password
	}

	update, err := dahuacore.UpdateDahuaCamera(models.DahuaCamera{
		ID:        camera.ID,
		Name:      form.Name,
		Address:   form.Address,
		Username:  form.Username,
		Password:  form.Password,
		Location:  location.Location,
		CreatedAt: camera.CreatedAt.Time,
		UpdatedAt: camera.UpdatedAt.Time,
	})
	if err != nil {
		return err
	}

	_, err = s.db.UpdateDahuaCamera(ctx, repo.UpdateDahuaCameraParams{
		ID:        update.ID,
		Name:      form.Name,
		Username:  update.Username,
		Password:  update.Password,
		Address:   update.Address,
		Location:  types.NewLocation(update.Location),
		UpdatedAt: types.NewTime(update.UpdatedAt),
	})
	if err != nil {
		return err
	}

	dbCamera, err := s.db.GetDahuaCamera(ctx, camera.ID)
	if err != nil {
		return err
	}
	s.bus.EventDahuaCameraUpdated(models.EventDahuaCameraUpdated{
		Camera: dbCamera.Convert(),
	})

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

func (s Server) DahuaEventsRulesCreatePOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Code string
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	err := dahua.CreateEventRule(ctx, s.db, repo.CreateDahuaEventRuleParams{
		Code:       form.Code,
		IgnoreDb:   true,
		IgnoreLive: true,
		IgnoreMqtt: true,
	})
	if err != nil {
		return err
	}

	if !isHTMX(c) {
		return c.Redirect(http.StatusSeeOther, "/dahua/events/rules")
	}

	return s.DahuaEventsRules(c)
}

func (s Server) DahuaEventsRulePOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Action string
		Rules  []struct {
			Selected bool
			ID       int64
			Code     string
			DB       bool
			Live     bool
			MQTT     bool
		}
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	switch form.Action {
	case "Update":
		for _, rule := range form.Rules {
			r, err := s.db.GetDahuaEventRule(ctx, rule.ID)
			if err != nil {
				if repo.IsNotFound(err) {
					continue
				}
				return err
			}

			err = dahua.UpdateEventRule(ctx, s.db, r, repo.UpdateDahuaEventRuleParams{
				Code:       rule.Code,
				IgnoreDb:   !rule.DB,
				IgnoreLive: !rule.Live,
				IgnoreMqtt: !rule.MQTT,
				ID:         rule.ID,
			})
			if err != nil {
				return err
			}
		}
	case "Delete":
		for _, rule := range form.Rules {
			if !rule.Selected {
				continue
			}

			rule, err := s.db.GetDahuaEventRule(ctx, rule.ID)
			if err != nil {
				if repo.IsNotFound(err) {
					continue
				}
				return err
			}

			if err := dahua.DeleteEventRule(ctx, s.db, rule); err != nil {
				return err
			}
		}
	}

	return s.DahuaEventsRules(c)
}

func (s Server) DahuaEventsRules(c echo.Context) error {
	ctx := c.Request().Context()

	rules, err := s.db.ListDahuaEventRule(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Rules": rules,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "dahua-events-rules", TemplateBlock{"htmx", data})
	}

	return c.Render(http.StatusOK, "dahua-events-rules", data)
}
