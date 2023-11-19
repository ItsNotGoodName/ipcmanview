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

func NewDahuaServer(dahuaStore *dahua.Store, dahuaPubSub PubSub) *DahuaServer {
	return &DahuaServer{
		dahuaStore:  dahuaStore,
		dahuaPubSub: dahuaPubSub,
	}
}

type DahuaServer struct {
	dahuaStore  *dahua.Store
	dahuaPubSub PubSub
}

func (s *DahuaServer) GET(c echo.Context) error {
	conns, err := s.dahuaStore.ConnList(c.Request().Context())
	if err != nil {
		return err
	}

	res := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		res = append(res, dahua.NewDahuaStatus(conn.Camera, conn.RPC.Conn))
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) POST(c echo.Context) error {
	var req map[string]models.DTODahuaCamera
	err := c.Bind(&req)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	ctx := c.Request().Context()
	for id, body := range req {
		camera, err := dahua.NewDahuaCamera(id, body)
		if err != nil {
			return echo.ErrBadRequest.WithInternal(err)
		}

		s.dahuaStore.ConnByCamera(ctx, camera)
	}

	return s.GET(c)
}

func (s *DahuaServer) IDPOST(c echo.Context) error {
	_, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct{}{})
}

func (s *DahuaServer) IDRPCPOST(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
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

func (s *DahuaServer) IDDetailGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetDahuaDetail(c.Request().Context(), conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) IDSoftwareGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetSoftwareVersion(c.Request().Context(), conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) IDLicensesGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetLicenseList(c.Request().Context(), conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) IDErrorGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	res := dahua.GetError(conn.RPC.Conn)

	return c.JSON(http.StatusOK, res)
}

func (s *DahuaServer) IDSnapshotGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
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

func (s *DahuaServer) IDEventsGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
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

			data := dahua.NewDahuaEvent(event, time.Now())

			if err := sendStream(c, stream, data); err != nil {
				return err
			}
		}
	} else {
		// Get events from the event worker

		ctx, cancel := context.WithCancel(c.Request().Context())
		defer cancel()

		dataC, err := s.dahuaPubSub.SubscribeDahuaEvents(ctx, []string{conn.Camera.ID})
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

func (s *DahuaServer) EventsGET(c echo.Context) error {
	ids := c.QueryParams()["ids"]

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	dataC, err := s.dahuaPubSub.SubscribeDahuaEvents(ctx, ids)
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

func (s *DahuaServer) IDFilesGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
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
				res, err := dahua.NewDahuaFiles(files, dahua.NewAffixSeed(conn.Camera.ID), conn.Camera.Location.Location)
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

func (s *DahuaServer) IDFilesPathGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	path := c.Param("*")

	req, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, dahuarpc.LoadFileURL(dahua.NewAddress(conn.Camera.Address), path), nil)
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

func (s *DahuaServer) IDAudioGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
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

func (s *DahuaServer) IDCoaxialStatusGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialStatus(c.Request().Context(), conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *DahuaServer) IDCoaxialCapsGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryInt(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialCaps(c.Request().Context(), conn.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *DahuaServer) IDPTZPresetPOST(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
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

func (s *DahuaServer) IDStorageGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	storage, err := dahua.GetStorage(c.Request().Context(), conn.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storage)
}

func (s *DahuaServer) IDUsersGET(c echo.Context) error {
	conn, err := useConn(c, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetUsers(c.Request().Context(), conn.RPC, conn.Camera.Location.Location)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
