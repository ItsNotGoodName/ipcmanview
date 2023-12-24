package webserver

import (
	"bytes"
	"context"
	"database/sql"
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

func (w Server) RegisterRoutes(e *echo.Echo) {
	e.GET("/", w.Index)
	e.GET("/dahua", w.Dahua)
	e.GET("/dahua/devices", w.DahuaDevices)
	e.GET("/dahua/devices/:id/update", w.DahuaDevicesUpdate)
	e.GET("/dahua/devices/create", w.DahuaDevicesCreate)
	e.GET("/dahua/devices/file-cursors", w.DahuaDevicesFileCursors)
	e.GET("/dahua/events", w.DahuaEvents)
	e.GET("/dahua/events/:id/data", w.DahuaEventsIDData)
	e.GET("/dahua/events/live", w.DahuaEventsLive)
	e.GET("/dahua/events/rules", w.DahuaEventsRules)
	e.GET("/dahua/events/stream", w.DahuaEventStream)
	e.GET("/dahua/files", w.DahuaFiles)
	e.GET("/dahua/snapshots", w.DahuaSnapshots)

	e.POST("/dahua/devices", w.DahuaDevicesPOST)
	e.POST("/dahua/devices/:id/update", w.DahuaDevicesUpdatePOST)
	e.POST("/dahua/devices/create", w.DahuaDevicesCreatePOST)
	e.POST("/dahua/devices/file-cursors", w.DahuaDevicesFileCursorsPOST)
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
		DeviceID  []int64
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
		DeviceID: params.DeviceID,
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

	devices, err := s.db.ListDahuaDevice(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Params":      params,
		"Devices":     devices,
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
		DeviceID  []int64
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
		DeviceID:  params.DeviceID,
		Start:     types.NewTime(timeRange.Start),
		End:       types.NewTime(timeRange.End),
		Ascending: params.Ascending,
		Local: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	devices, err := s.db.ListDahuaDevice(ctx)
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
		"Devices": devices,
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
		DeviceID []int64
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

	devices, err := s.db.ListDahuaDevice(ctx)
	if err != nil {
		return err
	}

	data := Data{
		"Params":      params,
		"Devices":     devices,
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
		DeviceID []int64
		Code     []string
		Action   []string
		Data     bool
	}{}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}

	sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaDeviceEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	w := useEventStream(c)
	buf := new(bytes.Buffer)

	for event := range eventsC {
		evt, ok := event.(models.EventDahuaDeviceEvent)
		if !ok ||
			evt.EventRule.IgnoreLive ||
			(len(params.DeviceID) > 0 && !slices.Contains(params.DeviceID, evt.Event.DeviceID)) ||
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

func (s Server) DahuaDevicesPOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Action  string
		Devices []struct {
			Selected bool
			ID       int64
		}
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	if form.Action == "Delete" {
		for _, device := range form.Devices {
			if !device.Selected {
				continue
			}
			if err := dahua.DeleteDevice(ctx, s.db, s.bus, device.ID); err != nil {
				return err
			}
		}
	}

	if isHTMX(c) {
		devices, err := s.db.ListDahuaDevice(c.Request().Context())
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-devices", TemplateBlock{"htmx-devices", Data{
			"Devices": devices,
		}})
	}

	return s.DahuaDevices(c)
}

func (s Server) DahuaDevices(c echo.Context) error {
	if isHTMX(c) {
		tables, err := useDahuaTables(c.Request().Context(), s.db, s.dahuaStore)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-devices", TemplateBlock{"htmx", tables})
	}

	devices, err := s.db.ListDahuaDevice(c.Request().Context())
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

	return c.Render(http.StatusOK, "dahua-devices", Data{
		"Devices":      devices,
		"FileCursors":  fileCursors,
		"EventWorkers": eventWorkers,
	})
}

func (s Server) DahuaDevicesCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-devices-create", Data{
		"Locations": dahua.Locations,
		"Features":  dahua.FeatureList,
	})
}

func (s Server) DahuaDevicesCreatePOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
		Features []string
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

	create, err := dahuacore.NewDahuaDevice(models.DahuaDevice{
		Name:     form.Name,
		Address:  form.Address,
		Username: form.Username,
		Password: form.Password,
		Location: location.Location,
		Feature:  dahua.FeatureFromStrings(form.Features),
	})

	err = dahua.CreateDevice(ctx, s.db, s.bus, repo.CreateDahuaDeviceParams{
		Name:      create.Name,
		Username:  create.Username,
		Password:  create.Password,
		Address:   create.Address,
		Location:  types.NewLocation(create.Location),
		Feature:   create.Feature,
		CreatedAt: types.NewTime(create.CreatedAt),
		UpdatedAt: types.NewTime(create.UpdatedAt),
	})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/dahua/devices")
}

func (s Server) DahuaDevicesFileCursorsPOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Action      string
		FileCursors []struct {
			Selected bool
			DeviceID int64
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

			err := dahua.ScanLockCreate(ctx, s.db, v.DeviceID)
			if err != nil {
				return err
			}

			err = dahua.ScanReset(ctx, s.db, v.DeviceID)

			dahua.ScanLockDelete(s.db, v.DeviceID)

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

			device, err := s.db.GetDahuaDevice(ctx, v.DeviceID)
			if err != nil {
				if repo.IsNotFound(err) {
					continue
				}
				return err
			}
			conn := s.dahuaStore.Conn(ctx, device.Convert().DahuaConn)

			if err := dahua.ScanLockCreate(ctx, s.db, v.DeviceID); err != nil {
				return err
			}
			go func(conn dahuacore.Conn) {
				ctx := context.Background()
				cancel := dahua.ScanLockHeartbeat(ctx, s.db, conn.Device.ID)
				defer cancel()

				err := dahua.Scan(ctx, s.db, conn.RPC, conn.Device, scanType)
				if err != nil {
					log.Err(err).Msg("Scan error")
				}
			}(conn)
		}
	}

	return s.DahuaDevicesFileCursors(c)
}

func (s Server) DahuaDevicesFileCursors(c echo.Context) error {
	if !isHTMX(c) {
		return c.Redirect(http.StatusSeeOther, "/dahua/devices#file-cursors")
	}

	fileCursors, err := s.db.ListDahuaFileCursor(c.Request().Context(), dahua.ScanLockStaleTime())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-devices", TemplateBlock{"htmx-file-cursors", Data{
		"FileCursors": fileCursors,
	}})
}

func (s Server) DahuaDevicesUpdate(c echo.Context) error {
	device, err := useDahuaDevice(c, s.db)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-devices-update", Data{
		"Locations": dahua.Locations,
		"Features":  dahua.FeatureList,
		"Device":    device,
	})
}

func (s Server) DahuaDevicesUpdatePOST(c echo.Context) error {
	ctx := c.Request().Context()

	device, err := useDahuaDevice(c, s.db)
	if err != nil {
		return err
	}

	var form struct {
		Name     string
		Address  string
		Username string
		Password string
		Location string
		Features []string
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}
	location, err := core.NewLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if form.Password == "" {
		form.Password = device.Password
	}

	update, err := dahuacore.UpdateDahuaDevice(models.DahuaDevice{
		ID:        device.ID,
		Name:      form.Name,
		Address:   form.Address,
		Username:  form.Username,
		Password:  form.Password,
		Location:  location.Location,
		Feature:   dahua.FeatureFromStrings(form.Features),
		CreatedAt: device.CreatedAt.Time,
		UpdatedAt: device.UpdatedAt.Time,
	})
	if err != nil {
		return err
	}

	err = dahua.UpdateDevice(ctx, s.db, s.bus, repo.UpdateDahuaDeviceParams{
		ID:        update.ID,
		Name:      form.Name,
		Username:  update.Username,
		Password:  update.Password,
		Address:   update.Address,
		Location:  types.NewLocation(update.Location),
		Feature:   update.Feature,
		UpdatedAt: types.NewTime(update.UpdatedAt),
	})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/dahua/devices")
}

func (s Server) DahuaSnapshots(c echo.Context) error {
	devices, err := s.db.ListDahuaDeviceByFeature(c.Request().Context(), models.DahuaFeatureCamera)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-snapshots", Data{
		"Devices": devices,
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
