package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	echo "github.com/labstack/echo/v4"
	"github.com/spf13/afero"
)

func (s *Server) DahuaAfero(route string) echo.HandlerFunc {
	return echo.WrapHandler(http.StripPrefix(route, http.FileServer(afero.NewHttpFs(s.dahuaFileFS))))
}

func (s *Server) DahuaDevices(c echo.Context) error {
	ctx := c.Request().Context()

	devices, err := s.db.DahuaListDevices(ctx)
	if err != nil {
		return err
	}
	clients := s.dahuaStore.ClientList(ctx, dahua.NewConns(devices))

	res := make([]models.DahuaStatus, 0, len(clients))
	for _, client := range clients {
		res = append(res, dahua.GetDahuaStatus(ctx, client.Conn, client.RPC))
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaDevicesIDRPCPOST(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
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

	res, err := dahuarpc.SendRaw[json.RawMessage](ctx, client.RPC, dahuarpc.New(req.Method).
		Params(req.Params).
		Object(req.Object))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaDevicesIDDetail(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetDahuaDetail(ctx, client.Conn.ID, client.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaDevicesIDSoftware(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetSoftwareVersion(ctx, client.Conn.ID, client.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaDevicesIDLicenses(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetLicenseList(ctx, client.Conn.ID, client.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaDevicesIDError(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res := dahua.GetError(ctx, client.RPC)

	return c.JSON(http.StatusOK, res)
}

func (s *Server) DahuaDevicesIDSnapshot(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	snapshot, err := dahuacgi.SnapshotGet(ctx, client.CGI, channel)
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

func (s *Server) DahuaDevicesIDEvents(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	direct, err := queryBoolOptional(c, "direct")
	if err != nil {
		return err
	}

	if direct {
		// Get events directly from the device

		manager, err := dahuacgi.EventManagerGet(ctx, client.CGI, 0)
		if err != nil {
			return err
		}
		reader := manager.Reader()

		stream := newStream(c)

		for {
			err := reader.Poll()
			if err != nil {
				return writeStreamError(c, stream, err)
			}

			event, err := reader.ReadEvent()
			if err != nil {
				return writeStreamError(c, stream, err)
			}

			// TODO: fix event json
			if err := writeStream(c, stream, event); err != nil {
				return err
			}
		}
	} else {
		// Get events from PubSub

		sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, event.DahuaEvent{})
		if err != nil {
			return err
		}
		defer sub.Close()

		stream := newStream(c)

		for e := range eventsC {
			evt, ok := e.(event.DahuaEvent)
			if !ok {
				continue
			}

			err := writeStream(c, stream, evt.Event)
			if err != nil {
				return writeStreamError(c, stream, err)
			}
		}

		if err := sub.Error(); err != nil {
			return writeStreamError(c, stream, err)
		}

		return nil
	}
}

func (s *Server) DahuaEvents(c echo.Context) error {
	ctx := c.Request().Context()

	ids, err := queryInts(c, "id")
	if err != nil {
		return err
	}

	sub, eventsC, err := s.pub.SubscribeChan(ctx, 10, event.DahuaEvent{})
	if err != nil {
		return err
	}
	defer sub.Close()

	stream := newStream(c)

	for e := range eventsC {
		evt, ok := e.(event.DahuaEvent)
		if !ok {
			continue
		}

		if len(ids) != 0 && !slices.Contains(ids, evt.Event.DeviceID) {
			continue
		}

		if err := writeStream(c, stream, evt.Event); err != nil {
			return writeStreamError(c, stream, err)
		}
	}

	if err := sub.Error(); err != nil {
		return writeStreamError(c, stream, err)
	}

	return nil
}

func (s *Server) DahuaDevicesIDFiles(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
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

	iterator := dahua.NewScannerPeriodIterator(timeRange)

	filesC := make(chan []mediafilefind.FindNextFileInfo)
	stream := newStream(c)

	for period, ok := iterator.Next(); ok; period, ok = iterator.Next() {
		cancel, errC := dahua.ScannerScan(ctx, client.RPC, period, client.Conn.Location, filesC)
		defer cancel()

	inner:
		for {
			select {
			case <-ctx.Done():
				return writeStreamError(c, stream, ctx.Err())
			case err := <-errC:
				if err != nil {
					return writeStreamError(c, stream, err)
				}
				break inner
			case files := <-filesC:
				res, err := dahua.NewDahuaFiles(client.Conn.ID, files, dahua.GetSeed(client.Conn), client.Conn.Location)
				if err != nil {
					return writeStreamError(c, stream, err)
				}

				if err := writeStream(c, stream, res); err != nil {
					return writeStreamError(c, stream, err)
				}
			}
		}
	}

	return nil
}

func (s *Server) DahuaDevicesIDFilesPath(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	filePath := c.Param("*")
	dbFile, err := s.db.DahuaGetFileByFilePath(ctx, repo.DahuaGetFileByFilePathParams{
		DeviceID: client.Conn.ID,
		FilePath: filePath,
	})
	if err != nil {
		if repo.IsNotFound(err) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	c.Response().Header().Set(echo.HeaderContentLength, strconv.FormatInt(dbFile.Length, 10))

	rd, err := func() (io.ReadCloser, error) {
		aferoFileFound := true
		aferoFile, err := s.db.DahuaGetAferoFileByFileID(ctx, sql.NullInt64{Int64: dbFile.ID, Valid: true})
		if repo.IsNotFound(err) {
			aferoFileFound = false
		} else if err != nil {
			return nil, err
		}

		if aferoFileFound {
			// File from cache
			rd, err := s.dahuaFileFS.Open(aferoFile.Name)
			if err != nil {
				if os.IsNotExist(err) {
					// File from device
					return dahua.FileLocalReadCloser(ctx, client, dbFile.FilePath)
				}

				return nil, err
			}

			return rd, nil
		}

		switch dbFile.Storage {
		case models.StorageLocal:
			// File from device
			return dahua.FileLocalReadCloser(ctx, client, dbFile.FilePath)
		case models.StorageFTP:
			// File from FTP
			return dahua.FileFTPReadCloser(ctx, s.db, dbFile.FilePath)
		case models.StorageSFTP:
			// File from SFTP
			return dahua.FileSFTPReadCloser(ctx, s.db, dbFile.FilePath)
		}

		return nil, echo.ErrInternalServerError.WithInternal(fmt.Errorf("storage not supported: %s", dbFile.FilePath))
	}()
	if err != nil {
		return err
	}
	defer rd.Close()

	_, err = io.Copy(c.Response().Writer, rd)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) DahuaDevicesIDAudio(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	audioStream, err := dahuacgi.AudioStreamGet(ctx, client.CGI, channel, dahuacgi.HTTPTypeSinglePart)
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

func (s *Server) DahuaDevicesIDCoaxialStatus(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialStatus(ctx, client.Conn.ID, client.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *Server) DahuaDevicesIDCoaxialCaps(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	channel, err := queryIntOptional(c, "channel")
	if err != nil {
		return err
	}

	status, err := dahua.GetCoaxialCaps(ctx, client.Conn.ID, client.RPC, channel)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}

func (s *Server) DahuaDevicesIDPTZPresetPOST(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
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

	err = dahua.SetPreset(ctx, client.PTZ, channel, index)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) DahuaDevicesIDStorage(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	storage, err := dahua.GetStorage(ctx, client.Conn.ID, client.RPC)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, storage)
}

func (s *Server) DahuaDevicesIDUsers(c echo.Context) error {
	ctx := c.Request().Context()

	client, err := useDahuaClient(c, s.db, s.dahuaStore)
	if err != nil {
		return err
	}

	res, err := dahua.GetUsers(ctx, client.Conn.ID, client.RPC, client.Conn.Location)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
