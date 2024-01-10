package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/files"
	"github.com/ItsNotGoodName/ipcmanview/internal/mediamtx"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/web/view"
	"github.com/ItsNotGoodName/ipcmanview/pkg/htmx"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func RegisterRenderer(e *echo.Echo) error {
	r, err := view.NewRenderer()
	if err != nil {
		return err
	}
	e.Renderer = r
	return nil
}

func (w Server) Register(e *echo.Echo) {
	e.DELETE("/dahua/events", w.DahuaEventsDELETE)
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
	e.GET("/dahua/files/:id/thumbnail", w.DahuaFilesIDThumbnail)
	e.GET("/dahua/snapshots", w.DahuaSnapshots)
	e.GET("/dahua/storage", w.DahuaStorage)
	e.GET("/dahua/storage/destinations/:id", w.DahuaStorageDestinationsID)
	e.GET("/dahua/storage/destinations/create", w.DahuaStorageDestinationsCreate)
	e.GET("/dahua/streams", w.DahuaStreams)
	e.PATCH("/dahua/devices/streams/:id", w.DahuaDevicesStreamsIDPATCH)
	e.POST("/dahua/devices", w.DahuaDevicesPOST)
	e.POST("/dahua/devices/:id/update", w.DahuaDevicesUpdatePOST)
	e.POST("/dahua/devices/create", w.DahuaDevicesCreatePOST)
	e.POST("/dahua/devices/file-cursors", w.DahuaDevicesFileCursorsPOST)
	e.POST("/dahua/events/rules", w.DahuaEventsRulePOST)
	e.POST("/dahua/events/rules/create", w.DahuaEventsRulesCreatePOST)
	e.POST("/dahua/files", w.DahuaFilesPOST)
	e.POST("/dahua/files/download", w.DahuaFilesDownloadPOST)
	e.POST("/dahua/storage/destinations", w.DahuaStorageDestinationsPOST)
	e.POST("/dahua/storage/destinations/:id", w.DahuaStorageDestinationsIDPOST)
	e.POST("/dahua/storage/destinations/:id/test", w.DahuaStorageDestinationsIDTestPOST)
	e.POST("/dahua/storage/destinations/create", w.DahuaStorageDestinationsCreatePOST)
}

type Server struct {
	db             repo.DB
	pub            pubsub.Pub
	bus            *core.Bus
	dahuaStore     *dahua.Store
	mediamtxConfig mediamtx.Config
	dahuaFileStore files.DahuaFileStore
}

func New(db repo.DB, pub pubsub.Pub, bus *core.Bus, dahuaStore *dahua.Store, dahuaFileStore files.DahuaFileStore, mediamtxConfig mediamtx.Config) Server {
	return Server{
		db:             db,
		pub:            pub,
		bus:            bus,
		dahuaStore:     dahuaStore,
		dahuaFileStore: dahuaFileStore,
		mediamtxConfig: mediamtxConfig,
	}
}

func (s Server) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (s Server) Dahua(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua", nil)
}

func (s Server) DahuaEventsDELETE(c echo.Context) error {
	ctx := c.Request().Context()

	err := s.db.DeleteDahuaEvent(ctx)
	if err != nil {
		return err
	}

	return s.DahuaEvents(c)
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

	data := view.Data{
		"Params":      params,
		"Devices":     devices,
		"Events":      events,
		"EventCodes":  eventCodes,
		"EventAction": eventActions,
	}

	if isHTMX(c) {
		htmx.SetReplaceURL(c.Response(), c.Request().URL.Path+"?"+api.EncodeQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-events", view.Block{Name: "htmx", Data: data})
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

	return c.Render(http.StatusOK, "dahua-events", view.Block{Name: "dahua-events-data-json", Data: data})
}

type DahuaFilesParams struct {
	DeviceID  []int64
	Type      []string
	Page      int64 `validate:"gt=0"`
	PerPage   int64 `validate:"gt=0"`
	Start     string
	End       string
	Ascending bool
}

func (params DahuaFilesParams) Filter() (repo.DahuaFileFilter, error) {
	timeRange, err := api.UseTimeRange(params.Start, params.End)
	if err != nil {
		return repo.DahuaFileFilter{}, err
	}

	return repo.DahuaFileFilter{
		Type:      params.Type,
		DeviceID:  params.DeviceID,
		Start:     types.NewTime(timeRange.Start),
		End:       types.NewTime(timeRange.End),
		Ascending: params.Ascending,
		Storage:   []models.Storage{},
	}, nil
}

func (s Server) DahuaFilesDownloadPOST(c echo.Context) error {
	ctx := c.Request().Context()

	params := DahuaFilesParams{}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}
	if err := api.ValidateStruct(params); err != nil {
		return err
	}

	filter, err := params.Filter()
	if err != nil {
		return err
	}

	for cursor := ""; ; {
		files, err := s.db.CursorListDahuaFile(ctx, repo.CursorListDahuaFileParams{
			PerPage:         100,
			Cursor:          cursor,
			DahuaFileFilter: filter,
		})
		if err != nil {
			return err
		}
		cursor = files.Cursor

		for _, dbFile := range files.Data {
			file := dbFile.Convert()

			exists, err := s.dahuaFileStore.Exists(ctx, file)
			if err != nil {
				return err
			}
			if exists {
				log.Info().Str("file-path", file.FilePath).Msg("Exists")
				continue
			}

			storage := core.StorageFromFilePath(file.FilePath)

			var rd io.ReadCloser
			switch storage {
			case models.StorageLocal:
				device, err := s.db.GetDahuaDevice(ctx, dbFile.DeviceID)
				if err != nil {
					return err
				}
				client := s.dahuaStore.Client(ctx, device.Convert().DahuaConn)

				rd, err = dahua.FileLocalReadCloser(ctx, client, file.FilePath)
				if err != nil {
					return err
				}
			case models.StorageFTP:
				rd, err = dahua.FileFTPReadCloser(ctx, s.db, file)
				if err != nil {
					return err
				}
			case models.StorageSFTP:
				rd, err = dahua.FileSFTPReadCloser(ctx, s.db, file)
				if err != nil {
					return err
				}
			default:
				log.Err(fmt.Errorf("invalid storage: %s", storage)).Send()
			}
			log.Info().Str("file-path", file.FilePath).Msg("Downloading...")
			err = s.dahuaFileStore.Save(ctx, file, rd)
			rd.Close()
			if err != nil {
				return err
			}
		}

		if !files.HasMore {
			break
		}
	}

	return c.NoContent(http.StatusOK)
}

func (s Server) DahuaFilesPOST(c echo.Context) error {
	s.bus.EventDahuaQuickScanQueue(models.EventDahuaQuickScanQueue{
		DeviceID: 0,
	})
	return s.DahuaFiles(c)
}

func (s Server) DahuaFiles(c echo.Context) error {
	ctx := c.Request().Context()

	params := DahuaFilesParams{
		Page:    1,
		PerPage: 10,
	}
	if err := api.DecodeQuery(c, &params); err != nil {
		return err
	}
	if err := api.ValidateStruct(params); err != nil {
		return err
	}

	filter, err := params.Filter()
	if err != nil {
		return err
	}

	files, err := s.db.ListDahuaFile(ctx, repo.ListDahuaFileParams{
		Page: pagination.Page{
			Page:    int(params.Page),
			PerPage: int(params.PerPage),
		},
		DahuaFileFilter: filter,
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

	data := view.Data{
		"Params":  params,
		"Files":   files,
		"Devices": devices,
		"Types":   types,
	}

	if isHTMX(c) {
		htmx.SetReplaceURL(c.Response(), c.Request().URL.Path+"?"+api.EncodeQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-files", view.Block{Name: "htmx", Data: data})
	}

	return c.Render(http.StatusOK, "dahua-files", data)
}

func (s Server) DahuaFilesIDThumbnail(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.ParamID(c)
	if err != nil {
		return err
	}

	file, err := s.db.GetDahuaFile(ctx, id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-files", view.Block{Name: "htmx-thumbnail",
		Data: file,
	})
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

	data := view.Data{
		"Params":      params,
		"Devices":     devices,
		"EventCodes":  eventCodes,
		"EventAction": eventActions,
	}

	if isHTMX(c) {
		htmx.SetReplaceURL(c.Response(), c.Request().URL.Path+"?"+api.EncodeQuery(params).Encode())
		return c.Render(http.StatusOK, "dahua-events-live", view.Block{Name: "htmx", Data: data})
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

	sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	w := useEventStream(c)
	buf := new(bytes.Buffer)

	for event := range eventsC {
		evt, ok := event.(models.EventDahuaEvent)
		if !ok ||
			evt.EventRule.IgnoreLive ||
			(len(params.DeviceID) > 0 && !slices.Contains(params.DeviceID, evt.Event.DeviceID)) ||
			(len(params.Code) > 0 && !slices.Contains(params.Code, evt.Event.Code)) ||
			(len(params.Action) > 0 && !slices.Contains(params.Action, evt.Event.Action)) {
			continue
		}

		type Event struct {
			DeviceName string
			models.DahuaEvent
		}

		if err := c.Echo().Renderer.Render(buf, "dahua-events-live", view.Block{Name: "event-row", Data: view.Data{
			"Event": Event{
				DeviceName: evt.DeviceName,
				DahuaEvent: evt.Event,
			},
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

		return c.Render(http.StatusOK, "dahua-devices", view.Block{Name: "htmx-devices", Data: view.Data{
			"Devices": devices,
		}})
	}

	return s.DahuaDevices(c)
}

func (s Server) DahuaDevices(c echo.Context) error {
	ctx := c.Request().Context()

	if isHTMX(c) {
		tables, err := useDahuaTables(ctx, s.db, s.dahuaStore)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "dahua-devices", view.Block{Name: "htmx", Data: tables})
	}

	devices, err := s.db.ListDahuaDevice(ctx)
	if err != nil {
		return err
	}

	fileCursors, err := s.db.ListDahuaFileCursor(ctx, dahua.ScanLockStaleTime())
	if err != nil {
		return err
	}

	eventWorkers, err := s.db.ListDahuaEventWorkerState(ctx)
	if err != nil {
		return err
	}

	streams, err := s.db.ListDahuaStream(ctx)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-devices", view.Data{
		"Devices":      devices,
		"FileCursors":  fileCursors,
		"EventWorkers": eventWorkers,
		"Streams":      streams,
	})
}

func (s Server) DahuaDevicesCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-devices-create", view.Data{
		"Locations": core.Locations,
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
	location, err := time.LoadLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	address, err := url.Parse(form.Address)
	if err != nil {
		return err
	}

	err = dahua.CreateDevice(ctx, s.db, s.bus, repo.CreateDahuaDeviceParams{
		Name:      form.Name,
		Username:  form.Username,
		Password:  form.Password,
		Address:   types.NewURL(address),
		Location:  types.NewLocation(location),
		Feature:   dahua.FeatureFromStrings(form.Features),
		CreatedAt: types.NewTime(time.Now()),
		UpdatedAt: types.NewTime(time.Now()),
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

			err := dahua.ScanLockCreateTry(ctx, s.db, v.DeviceID)
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
	case "Quick":
		for _, v := range form.FileCursors {
			if !v.Selected {
				continue
			}

			s.bus.EventDahuaQuickScanQueue(models.EventDahuaQuickScanQueue{
				DeviceID: v.DeviceID,
			})
		}
	case "Full", "Reverse":
		scanType := models.DahuaScanTypeFull
		if form.Action == "Full" {
			scanType = models.DahuaScanTypeFull
		} else if form.Action == "Reverse" {
			scanType = models.DahuaScanTypeReverse
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
			client := s.dahuaStore.Client(ctx, device.Convert().DahuaConn)

			if err := dahua.ScanLockCreateTry(ctx, s.db, v.DeviceID); err != nil {
				return err
			}
			go func(conn dahua.Client) {
				ctx := context.Background()
				cancel := dahua.ScanLockHeartbeat(ctx, s.db, conn.Conn.ID)
				defer cancel()

				err := dahua.Scan(ctx, s.db, conn.RPC, conn.Conn, scanType)
				if err != nil {
					log.Err(err).Msg("Failed to scan")
				}
			}(client)
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

	return c.Render(http.StatusOK, "dahua-devices", view.Block{Name: "htmx-file-cursors", Data: view.Data{
		"FileCursors": fileCursors,
	}})
}

func (s Server) DahuaDevicesUpdate(c echo.Context) error {
	device, err := useDahuaDevice(c, s.db)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-devices-update", view.Data{
		"Locations": core.Locations,
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
	location, err := time.LoadLocation(form.Location)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	address, err := url.Parse(form.Address)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	if form.Password == "" {
		form.Password = device.Password
	}

	err = dahua.UpdateDevice(ctx, s.db, s.bus, device.Convert().DahuaDevice, repo.UpdateDahuaDeviceParams{
		ID:        device.ID,
		Name:      form.Name,
		Username:  form.Username,
		Password:  form.Password,
		Address:   types.NewURL(address),
		Location:  types.NewLocation(location),
		Feature:   dahua.FeatureFromStrings(form.Features),
		UpdatedAt: types.NewTime(time.Now()),
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

	return c.Render(http.StatusOK, "dahua-snapshots", view.Data{
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

	rules, err := s.db.ListDahuaEventRule(ctx)
	if err != nil {
		return err
	}

	data := view.Data{
		"Rules": rules,
	}

	return c.Render(http.StatusOK, "dahua-events-rules", view.Block{Name: "htmx-create", Data: data})
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

	data := view.Data{
		"Rules": rules,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "dahua-events-rules", view.Block{Name: "htmx", Data: data})
	}

	return c.Render(http.StatusOK, "dahua-events-rules", data)
}

func (s Server) DahuaStreams(c echo.Context) error {
	ctx := c.Request().Context()

	type ViewCamera struct {
		models.DahuaDevice
		SelectedStream *models.DahuaStream
		Streams        []models.DahuaStream
	}

	if isHTMX(c) {
		var params struct {
			ID int64
		}
		err := api.DecodeQuery(c, &params)
		if err != nil {
			return err
		}

		dbStream, err := s.db.GetDahuaStream(ctx, params.ID)
		if err != nil {
			return err
		}

		dbCamera, err := s.db.GetDahuaDevice(ctx, dbStream.DeviceID)
		if err != nil {
			return err
		}

		dbStreams, err := s.db.ListDahuaStreamByDevice(ctx, dbStream.DeviceID)
		if err != nil {
			return err
		}
		streams := make([]models.DahuaStream, 0, len(dbStreams))
		for _, dbStream := range dbStreams {
			streams = append(streams, dbStream.Convert(s.mediamtxConfig.DahuaEmbedURL(dbStream)))
		}

		stream := dbStream.Convert(s.mediamtxConfig.DahuaEmbedURL(dbStream))

		viewDevice := ViewCamera{
			DahuaDevice:    dbCamera.Convert().DahuaDevice,
			SelectedStream: &stream,
			Streams:        streams,
		}

		return c.Render(http.StatusOK, "dahua-streams", view.Block{Name: "htmx", Data: viewDevice})
	}

	dbCameras, err := s.db.ListDahuaDeviceByFeature(ctx, models.DahuaFeatureCamera)
	if err != nil {
		return err
	}

	viewDevices := make([]ViewCamera, 0, len(dbCameras))
	for _, camera := range dbCameras {
		dbStreams, err := s.db.ListDahuaStreamByDevice(ctx, camera.ID)
		if err != nil {
			return err
		}

		streams := make([]models.DahuaStream, 0, len(dbStreams))
		for _, dbStream := range dbStreams {
			streams = append(streams, dbStream.Convert(s.mediamtxConfig.DahuaEmbedURL(dbStream)))
		}

		var selectedStream *models.DahuaStream
		if len(streams) > 0 {
			selectedStream = &streams[0]
		}

		viewDevices = append(viewDevices, ViewCamera{
			DahuaDevice:    camera.Convert().DahuaDevice,
			SelectedStream: selectedStream,
			Streams:        streams,
		})
	}

	return c.Render(http.StatusOK, "dahua-streams", view.Data{
		"Devices": viewDevices,
	})
}

func (s Server) DahuaStorage(c echo.Context) error {
	ctx := c.Request().Context()

	storageDestinations, err := s.db.ListDahuaStorageDestination(ctx)
	if err != nil {
		return err
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "dahua-storage", view.Block{Name: "htmx-storage-destination", Data: view.Data{
			"StorageDestinations": storageDestinations,
		}})
	}

	return c.Render(http.StatusOK, "dahua-storage", view.Data{
		"StorageDestinations": storageDestinations,
	})
}

func (s Server) DahuaStorageDestinationsPOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Action              string
		StorageDestinations []struct {
			Selected bool
			ID       int64
		}
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	switch form.Action {
	case "Delete":
		for _, v := range form.StorageDestinations {
			if !v.Selected {
				continue
			}

			err := dahua.DeleteStorageDestination(ctx, s.db, v.ID)
			if err != nil {
				return err
			}
		}
	}

	return s.DahuaStorage(c)
}

func (s Server) DahuaStorageDestinationsID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.ParamID(c)
	if err != nil {
		return err
	}

	storageDestination, err := s.db.GetDahuaStorageDestination(ctx, id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "dahua-storage-destinations-id", view.Data{
		"StorageDestination": storageDestination,
		"Storage":            dahua.Storage,
	})
}

func (s Server) DahuaStorageDestinationsCreate(c echo.Context) error {
	return c.Render(http.StatusOK, "dahua-storage-destinations-create", view.Data{
		"Storage": dahua.Storage,
	})
}

func (s Server) DahuaStorageDestinationsCreatePOST(c echo.Context) error {
	ctx := c.Request().Context()

	var form struct {
		Name            string
		Storage         string
		ServerAddress   string
		Port            int64
		Username        string
		Password        string
		RemoteDirectory string
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}
	storage, err := dahua.ParseStorage(form.Storage)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	_, err = dahua.CreateStorageDestination(ctx, s.db, models.DahuaStorageDestination{
		Name:            form.Name,
		Storage:         storage,
		ServerAddress:   form.ServerAddress,
		Port:            form.Port,
		Username:        form.Username,
		Password:        form.Password,
		RemoteDirectory: form.RemoteDirectory,
	})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/dahua/storage")
}

func (s Server) DahuaStorageDestinationsIDPOST(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.ParamID(c)
	if err != nil {
		return err
	}
	storageDestination, err := s.db.GetDahuaStorageDestination(ctx, id)
	if err != nil {
		if repo.IsNotFound(err) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	var form struct {
		Name            string
		Storage         string
		ServerAddress   string
		Port            int64
		Username        string
		Password        string
		RemoteDirectory string
		UpdatePassword  bool
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}
	if !form.UpdatePassword {
		form.Password = storageDestination.Password
	}
	storage, err := dahua.ParseStorage(form.Storage)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	storageDestination.Name = form.Name
	storageDestination.Storage = storage
	storageDestination.ServerAddress = form.ServerAddress
	storageDestination.Port = form.Port
	storageDestination.Username = form.Username
	storageDestination.Password = form.Password
	storageDestination.RemoteDirectory = form.RemoteDirectory

	_, err = dahua.UpdateStorageDestination(ctx, s.db, storageDestination.Convert())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/dahua/storage")
}

func (s Server) DahuaStorageDestinationsIDTestPOST(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.ParamID(c)
	if err != nil {
		return err
	}
	storageDestination, err := s.db.GetDahuaStorageDestination(ctx, id)
	if err != nil {
		if repo.IsNotFound(err) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	err = dahua.TestStorageDestination(ctx, storageDestination.Convert())
	if err != nil {
		htmx.NewEvent("toast-error", err.Error()).SetTrigger(c.Response())
		return c.Render(http.StatusOK, "dahua-storage", view.Block{Name: "htmx-storage-destination-test", Data: view.Data{
			"OK": false,
			"ID": storageDestination.ID,
		}})
	} else {
		htmx.NewEvent("toast", "OK").SetTrigger(c.Response())
		return c.Render(http.StatusOK, "dahua-storage", view.Block{Name: "htmx-storage-destination-test", Data: view.Data{
			"OK": true,
			"ID": storageDestination.ID,
		}})
	}
}

func (s Server) DahuaDevicesStreamsIDPATCH(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.ParamID(c)
	if err != nil {
		return err
	}
	stream, err := s.db.GetDahuaStream(ctx, id)
	if err != nil {
		return err
	}

	var form struct {
		Name         *string
		MediamtxPath *string
	}
	if err := api.ParseForm(c, &form); err != nil {
		return err
	}

	stream, err = dahua.UpdateStream(ctx, s.db, stream, repo.UpdateDahuaStreamParams{
		Name:         repo.Coalasce(form.Name, &stream.Name),
		MediamtxPath: repo.Coalasce(form.MediamtxPath, &stream.MediamtxPath),
		ID:           stream.ID,
	})
	if err != nil {
		return err
	}

	if form.Name != nil {
		return c.Render(http.StatusOK, "dahua-devices", view.Block{Name: "htmx-streams-name", Data: stream})
	} else if form.MediamtxPath != nil {
		return c.Render(http.StatusOK, "dahua-devices", view.Block{Name: "htmx-streams-mediamtx-path", Data: stream})
	} else {
		return c.NoContent(http.StatusNoContent)
	}
}
