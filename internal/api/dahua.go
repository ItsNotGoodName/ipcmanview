package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	echo "github.com/labstack/echo/v4"
)

type DahuaRepo interface {
	GetConn(ctx context.Context, id int64) (models.DahuaConn, bool, error)
	ListConn(ctx context.Context) ([]models.DahuaConn, error)
	GetFileByFilePath(ctx context.Context, deviceID int64, filePath string) (models.DahuaFile, error)
}

type DahuaFileCache interface {
	// Save(ctx context.Context, file models.DahuaFile, r io.ReadCloser) error
	Exists(ctx context.Context, file models.DahuaFile) (bool, error)
	Get(ctx context.Context, file models.DahuaFile) (io.ReadCloser, error)
}

func NewDahuaServer(
	pub pubsub.Pub,
	dahuaStore *dahuacore.Store,
	dahuaRepo DahuaRepo,
	dahuaFileCache DahuaFileCache,
) *DahuaServer {
	return &DahuaServer{
		pub:            pub,
		dahuaStore:     dahuaStore,
		dahuaRepo:      dahuaRepo,
		dahuaFileCache: dahuaFileCache,
	}
}

type DahuaServer struct {
	pub            pubsub.Pub
	dahuaStore     *dahuacore.Store
	dahuaRepo      DahuaRepo
	dahuaFileCache DahuaFileCache
}

func (s *DahuaServer) Register(e *echo.Echo) {
	e.GET("/v1/dahua", s.Dahua)
	e.GET("/v1/dahua-events", s.DahuaEvents)
	e.GET("/v1/dahua/:id/audio", s.DahuaIDAudio)
	e.GET("/v1/dahua/:id/coaxial/caps", s.DahuaIDCoaxialCaps)
	e.GET("/v1/dahua/:id/coaxial/status", s.DahuaIDCoaxialStatus)
	e.GET("/v1/dahua/:id/detail", s.DahuaIDDetail)
	e.GET("/v1/dahua/:id/error", s.DahuaIDError)
	e.GET("/v1/dahua/:id/events", s.DahuaIDEvents)
	e.GET("/v1/dahua/:id/files", s.DahuaIDFiles)
	e.GET("/v1/dahua/:id/files/*", s.DahuaIDFilesPath)
	e.GET("/v1/dahua/:id/licenses", s.DahuaIDLicenses)
	e.GET("/v1/dahua/:id/snapshot", s.DahuaIDSnapshot)
	e.GET("/v1/dahua/:id/software", s.DahuaIDSoftware)
	e.GET("/v1/dahua/:id/storage", s.DahuaIDStorage)
	e.GET("/v1/dahua/:id/users", s.DahuaIDUsers)

	e.POST("/v1/dahua/:id/ptz/preset", s.DahuaIDPTZPresetPOST)
	e.POST("/v1/dahua/:id/rpc", s.DahuaIDRPCPOST)
}

func useDahuaConn(c echo.Context, repo DahuaRepo, store *dahuacore.Store) (dahuacore.Client, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return dahuacore.Client{}, echo.ErrBadRequest.WithInternal(err)
	}

	ctx := c.Request().Context()

	device, found, err := repo.GetConn(ctx, id)
	if err != nil {
		return dahuacore.Client{}, err
	}
	if !found {
		return dahuacore.Client{}, echo.ErrNotFound.WithInternal(err)
	}

	client := store.Conn(ctx, device)

	return client, nil
}

func (s *DahuaServer) Dahua(c echo.Context) error {
	ctx := c.Request().Context()

	devices, err := s.dahuaRepo.ListConn(ctx)
	if err != nil {
		return err
	}

	conns := s.dahuaStore.ConnList(ctx, devices)

	res := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		res = append(res, dahuacore.GetDahuaStatus(conn.Conn, conn.RPC))
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) DahuaIDRPCPOST(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	var req struct {
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
		Object int64           `json:"object"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := dahuarpc.SendRaw[json.RawMessage](ctx, conn.RPC, dahuarpc.New(req.Method).
		Params(req.Params).
		Object(req.Object))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) DahuaIDDetail(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahuacore.GetDahuaDetail(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) DahuaIDSoftware(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahuacore.GetSoftwareVersion(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) DahuaIDLicenses(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahuacore.GetLicenseList(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) DahuaIDError(c echo.Context) error {
	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	res := dahuacore.GetError(conn.RPC)

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) DahuaIDSnapshot(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	snapshot, err := dahuacgi.SnapshotGet(ctx, conn.CGI, channel)
	if err != nil {
		return err
	}
	defer snapshot.Close()

	c.Response().Header().Set(echo.HeaderContentLength, snapshot.ContentLength)
	c.Response().Header().Set(echo.HeaderCacheControl, "no-store")

	_, err = io.Copy(c.Response().Writer, snapshot)
	if err != nil {
		return err
	}

	return nil
}

func (s *DahuaServer) DahuaIDEvents(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	direct, err := queryBoolOptional(c, "direct")
	if err != nil {
		return err
	}

	if direct {
		// Get events directly from the device

		manager, err := dahuacgi.EventManagerGet(ctx, conn.CGI, 0)
		if err != nil {
			return err
		}
		reader := manager.Reader()

		stream := useStream(c)

		for {
			err := reader.Poll()
			if err != nil {
				return sendStreamError(c, stream, err)
			}

			event, err := reader.ReadEvent()
			if err != nil {
				return sendStreamError(c, stream, err)
			}

			data := dahuacore.NewDahuaEvent(conn.Conn.ID, event)

			if err := sendStream(c, stream, data); err != nil {
				return err
			}
		}
	} else {
		// Get events from PubSub

		sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaDeviceEvent{})
		if err != nil {
			return err
		}
		defer sub.Close()

		stream := useStream(c)

		for event := range eventsC {
			evt, ok := event.(models.EventDahuaDeviceEvent)
			if !ok {
				continue
			}

			err := sendStream(c, stream, evt.Event)
			if err != nil {
				return sendStreamError(c, stream, err)
			}
		}

		if err := sub.Error(); err != nil {
			return sendStreamError(c, stream, err)
		}

		return nil
	}
}

func (s *DahuaServer) DahuaEvents(c echo.Context) error {
	ctx := c.Request().Context()

	ids, err := queryInts(c, "id")
	if err != nil {
		return err
	}

	sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, models.EventDahuaDeviceEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	stream := useStream(c)

	for event := range eventsC {
		evt, ok := event.(models.EventDahuaDeviceEvent)
		if !ok {
			continue
		}

		if len(ids) != 0 && !slices.Contains(ids, evt.Event.DeviceID) {
			continue
		}

		if err := sendStream(c, stream, evt.Event); err != nil {
			return sendStreamError(c, stream, err)
		}
	}

	if err := sub.Error(); err != nil {
		return sendStreamError(c, stream, err)
	}

	return nil
}

func (s *DahuaServer) DahuaIDFiles(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	var form struct {
		Start string
		End   string
	}
	if err := DecodeQuery(c, &form); err != nil {
		return err
	}

	timeRange, err := UseTimeRange(form.Start, form.End)
	if err != nil {
		return err
	}

	iterator := dahuacore.NewScanPeriodIterator(timeRange)

	filesC := make(chan []mediafilefind.FindNextFileInfo)
	stream := useStream(c)

	for period, ok := iterator.Next(); ok; period, ok = iterator.Next() {
		cancel, errC := dahuacore.Scan(ctx, conn.RPC, period, conn.Conn.Location, filesC)
		defer cancel()

	inner:
		for {
			select {
			case <-ctx.Done():
				return sendStreamError(c, stream, ctx.Err())
			case err := <-errC:
				if err != nil {
					return sendStreamError(c, stream, err)
				}
				break inner
			case files := <-filesC:
				res, err := dahuacore.NewDahuaFiles(conn.Conn.ID, files, dahuacore.GetSeed(conn.Conn), conn.Conn.Location)
				if err != nil {
					return sendStreamError(c, stream, err)
				}

				if err := sendStream(c, stream, res); err != nil {
					return sendStreamError(c, stream, err)
				}
			}
		}
	}

	return nil
}

func (s *DahuaServer) DahuaIDFilesPath(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	filePath := c.Param("*")

	dahuaFile, err := s.dahuaRepo.GetFileByFilePath(ctx, conn.Conn.ID, filePath)
	if err != nil {
		return err
	}

	var rd io.ReadCloser
	if exists, err := s.dahuaFileCache.Exists(ctx, dahuaFile); err != nil {
		return err
	} else if !exists {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dahuarpc.LoadFileURL(dahuacore.NewHTTPAddress(conn.Conn.Address), filePath), nil)
		if err != nil {
			return err
		}

		session, err := conn.RPC.Session(ctx)
		if err != nil {
			return err
		}

		req.Header.Add("Cookie", dahuarpc.Cookie(session))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		rd = resp.Body
	} else {
		rd, err = s.dahuaFileCache.Get(ctx, dahuaFile)
		if err != nil {
			return err
		}
	}
	defer rd.Close()

	_, err = io.Copy(c.Response().Writer, rd)
	if err != nil {
		return err
	}

	return nil
}

func (s *DahuaServer) DahuaIDAudio(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	audioStream, err := dahuacgi.AudioStreamGet(ctx, conn.CGI, channel, dahuacgi.HTTPTypeSinglePart)
	if err != nil {
		return err
	}

	c.Request().Header.Add("ContentType", audioStream.ContentType)

	_, err = io.Copy(c.Response().Writer, audioStream)
	if err != nil {
		return err
	}

	return nil
}

func (s *DahuaServer) DahuaIDCoaxialStatus(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahuacore.GetCoaxialStatus(ctx, conn.Conn.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *DahuaServer) DahuaIDCoaxialCaps(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahuacore.GetCoaxialCaps(ctx, conn.Conn.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *DahuaServer) DahuaIDPTZPresetPOST(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	index, err := queryIntOptional(c, "index")
	if err != nil {
		return err
	}

	err = dahuacore.SetPreset(ctx, conn.PTZ, channel, index)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *DahuaServer) DahuaIDStorage(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	storage, err := dahuacore.GetStorage(ctx, conn.Conn.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storage)
}

func (s *DahuaServer) DahuaIDUsers(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := useDahuaConn(c, s.dahuaRepo, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahuacore.GetUsers(ctx, conn.Conn.ID, conn.RPC, conn.Conn.Location)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
