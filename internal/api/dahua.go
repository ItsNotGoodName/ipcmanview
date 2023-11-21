package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	echo "github.com/labstack/echo/v4"
)

type DahuaStore interface {
	ConnListByCameras(ctx context.Context, cameras ...models.DahuaCamera) ([]dahua.Conn, error)
	ConnByID(ctx context.Context, id string) (dahua.Conn, error)
}

func RegisterDahuaRoutes(e *echo.Echo, s *DahuaServer) {
	e.GET("/v1/dahua", s.GET)
	e.GET("/v1/dahua-events", s.GETEvents)
	e.GET("/v1/dahua/:id/audio", s.GETIDAudio)
	e.GET("/v1/dahua/:id/coaxial/caps", s.GETIDCoaxialCaps)
	e.GET("/v1/dahua/:id/coaxial/status", s.GETIDCoaxialStatus)
	e.GET("/v1/dahua/:id/detail", s.GETIDDetail)
	e.GET("/v1/dahua/:id/error", s.GETIDError)
	e.GET("/v1/dahua/:id/events", s.GETIDEvents)
	e.GET("/v1/dahua/:id/files", s.GETIDFiles)
	e.GET("/v1/dahua/:id/files/*", s.GETIDFilesPath)
	e.GET("/v1/dahua/:id/licenses", s.GETIDLicenses)
	e.GET("/v1/dahua/:id/snapshot", s.GETIDSnapshot)
	e.GET("/v1/dahua/:id/software", s.GETIDSoftware)
	e.GET("/v1/dahua/:id/storage", s.GETIDStorage)
	e.GET("/v1/dahua/:id/users", s.GETIDUsers)

	e.POST("/v1/dahua", s.POST)
	e.POST("/v1/dahua/:id", s.POSTID)
	e.POST("/v1/dahua/:id/ptz/preset", s.POSTIDPTZPreset)
	e.POST("/v1/dahua/:id/rpc", s.POSTIDRPC)
}

func useDahuaConn(c echo.Context, store DahuaStore) (dahua.Conn, error) {
	id := c.Param("id")

	client, err := store.ConnByID(c.Request().Context(), id)
	if err != nil {
		return dahua.Conn{}, echo.ErrNotFound.WithInternal(err)
	}

	return client, nil
}

func NewDahuaServer(dahuaStore DahuaStore, dahuaPubSub PubSub) *DahuaServer {
	return &DahuaServer{
		store:  dahuaStore,
		pubSub: dahuaPubSub,
	}
}

type DahuaServer struct {
	store  DahuaStore
	pubSub PubSub
}

func (s *DahuaServer) GET(c echo.Context) error {
	conns, err := s.store.ConnListByCameras(c.Request().Context())
	if err != nil {
		return err
	}

	res := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		res = append(res, dahua.GetDahuaStatus(conn.Camera, conn.RPC.Conn))
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) POST(c echo.Context) error {
	var req map[string]models.DTODahuaCamera
	err := c.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	cameras := make([]models.DahuaCamera, 0, len(req))
	for id, body := range req {
		camera, err := dahua.NewDahuaCamera(id, body)
		if err != nil {
			return echo.ErrBadRequest.WithInternal(err)
		}

		cameras = append(cameras, camera)
	}

	_, err = s.store.ConnListByCameras(c.Request().Context(), cameras...)
	if err != nil {
		return err
	}

	return s.GET(c)
}

func (s *DahuaServer) POSTID(c echo.Context) error {
	_, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct{}{})
}

func (s *DahuaServer) POSTIDRPC(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	var req struct {
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
		Object int64           `json:"object"`
		Seq    int             `json:"seq"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	ctx := c.Request().Context()

	rpc, err := conn.RPC.RPC(ctx)
	if err != nil {
		return err
	}

	res, err := dahuarpc.SendRaw[json.RawMessage](ctx, rpc.
		Method(req.Method).
		Params(req.Params).
		Object(req.Object).
		Seq(req.Seq))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) GETIDDetail(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	res, err := dahua.GetDahuaDetail(c.Request().Context(), conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) GETIDSoftware(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	res, err := dahua.GetSoftwareVersion(c.Request().Context(), conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) GETIDLicenses(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	res, err := dahua.GetLicenseList(c.Request().Context(), conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) GETIDError(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	res := dahua.GetError(conn.RPC.Conn)

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) GETIDSnapshot(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	snapshot, err := dahuacgi.SnapshotGet(c.Request().Context(), conn.CGI, channel)
	if err != nil {
		return err
	}
	defer snapshot.Close()

	c.Response().Header().Set(echo.HeaderContentLength, snapshot.ContentLength)

	_, err = io.Copy(c.Response().Writer, snapshot)
	if err != nil {
		return err
	}

	return nil
}

func (s *DahuaServer) GETIDEvents(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	direct, err := queryBool(c, "direct")
	if err != nil {
		return err
	}

	if direct {
		// Get events directly from the camera

		manager, err := dahuacgi.EventManagerGet(c.Request().Context(), conn.CGI, 0)
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

			data := dahua.NewDahuaEvent(conn.Camera.ID, event, time.Now())

			if err := sendStream(c, stream, data); err != nil {
				return err
			}
		}
	} else {
		// Get events from the event worker

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		dataC, err := s.pubSub.SubscribeDahuaEvents(ctx, []string{conn.Camera.ID})
		if err != nil {
			return err
		}

		stream := useStream(c)

		for {
			select {
			case <-ctx.Done():
				return sendStreamError(c, stream, ctx.Err())
			case data, ok := <-dataC:
				if !ok {
					return sendStreamError(c, stream, ErrSubscriptionClosed)
				}
				if err := sendStream(c, stream, data.Event); err != nil {
					return err
				}
			}
		}
	}
}

func (s *DahuaServer) GETEvents(c echo.Context) error {
	ids := c.QueryParams()["ids"]

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	dataC, err := s.pubSub.SubscribeDahuaEvents(ctx, ids)
	if err != nil {
		return err
	}

	stream := useStream(c)

	for {
		select {
		case <-ctx.Done():
			return sendStreamError(c, stream, ctx.Err())
		case data, ok := <-dataC:
			if !ok {
				return sendStreamError(c, stream, ErrSubscriptionClosed)
			}
			if err := sendStream(c, stream, data); err != nil {
				return err
			}
		}
	}
}

func (s *DahuaServer) GETIDFiles(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	scanRange, err := queryDahuaScanRange(c)
	if err != nil {
		return err
	}
	iter := dahua.NewScanPeriodIterator(scanRange)

	filesC := make(chan []mediafilefind.FindNextFileInfo)
	stream := useStream(c)
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	for period, ok := iter.Next(); ok; period, ok = iter.Next() {
		errC := dahua.Scan(ctx, conn.RPC, period, conn.Camera.Location.Location, filesC)

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
				res, err := dahua.NewDahuaFiles(conn.Camera.ID, files, dahua.GetSeed(conn.Camera), conn.Camera.Location.Location)
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

func (s *DahuaServer) GETIDFilesPath(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	path := c.Param("*")

	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, dahuarpc.LoadFileURL(dahua.NewHTTPAddress(conn.Camera.Address), path), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Cookie", dahuarpc.Cookie(conn.RPC.Data().Session))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *DahuaServer) GETIDAudio(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	audioStream, err := dahuacgi.AudioStreamGet(c.Request().Context(), conn.CGI, channel, dahuacgi.HTTPTypeSinglePart)
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

func (s *DahuaServer) GETIDCoaxialStatus(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialStatus(c.Request().Context(), conn.Camera.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *DahuaServer) GETIDCoaxialCaps(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialCaps(c.Request().Context(), conn.Camera.ID, conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *DahuaServer) POSTIDPTZPreset(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	index, err := queryInt(c, "index")
	if err != nil {
		return err
	}

	err = dahua.SetPreset(c.Request().Context(), conn.PTZ, channel, index)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *DahuaServer) GETIDStorage(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	storage, err := dahua.GetStorage(c.Request().Context(), conn.Camera.ID, conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storage)
}

func (s *DahuaServer) GETIDUsers(c echo.Context) error {
	conn, err := useDahuaConn(c, s.store)
	if err != nil {
		return err
	}

	res, err := dahua.GetUsers(c.Request().Context(), conn.Camera.ID, conn.RPC, conn.Camera.Location.Location)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
