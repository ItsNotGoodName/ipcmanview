package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/labstack/echo/v4"
)

func paramID(c echo.Context) (int64, error) {
	number, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func useDahuaPermissions(c echo.Context, db repo.DB) (models.DahuaDevicePermissions, error) {
	ctx := c.Request().Context()

	session, ok := auth.UseSession(ctx)
	if !ok {
		return nil, echo.ErrUnauthorized
	}

	permissions, err := db.DahuaListDahuaDevicePermissions(ctx, repo.DahuaDevicePermissionParams{
		UserID: session.UserID,
		Level:  models.DahuaPermissionLevelUser,
	})
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func useDahuaClient(c echo.Context, db repo.DB, store *dahua.Store) (dahua.Client, models.DahuaPermissionLevel, error) {
	ctx := c.Request().Context()

	id, err := paramID(c)
	if err != nil {
		return dahua.Client{}, 0, err
	}

	permissions, err := useDahuaPermissions(c, db)

	permission, ok := permissions.Get(id)
	if !ok {
		return dahua.Client{}, 0, echo.ErrNotFound
	}

	conn, err := dahua.GetFatDevice(ctx, db, permissions, repo.DahuaFatDeviceParams{IDs: []int64{id}})
	if err != nil {
		if repo.IsNotFound(err) {
			return dahua.Client{}, 0, echo.ErrNotFound.WithInternal(err)
		}
		return dahua.Client{}, 0, err
	}

	return store.Client(ctx, dahua.NewConn(conn)), permission.Level, nil
}

// ---------- Stream

func newStream(c echo.Context) *json.Encoder {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response())
}

type StreamPayload struct {
	Data    any     `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	OK      bool    `json:"ok"`
}

func writeStreamError(c echo.Context, enc *json.Encoder, err error) error {
	str := err.Error()
	if encodeErr := enc.Encode(StreamPayload{
		OK:      false,
		Message: &str,
	}); encodeErr != nil {
		return errors.Join(encodeErr, err)
	}

	c.Response().Flush()

	return err
}

func writeStream(c echo.Context, enc *json.Encoder, data any) error {
	err := enc.Encode(StreamPayload{
		OK:   true,
		Data: data,
	})
	if err != nil {
		return writeStreamError(c, enc, err)
	}

	c.Response().Flush()

	return nil
}

// ---------- Queries

func queryTimeRange(c echo.Context) (models.TimeRange, error) {
	var query struct {
		Start string
		End   string
	}
	if err := c.Bind(&query); err != nil {
		return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
	}
	start, end := query.Start, query.End

	var startTime, endTime time.Time
	if start != "" {
		var err error
		startTime, err = time.ParseInLocation("2006-01-02T15:04", start, time.Local)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	}

	if end != "" {
		var err error
		endTime, err = time.ParseInLocation("2006-01-02T15:04", end, time.Local)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	} else if start != "" {
		endTime = time.Now()
	}

	r, err := core.NewTimeRange(startTime, endTime)
	if err != nil {
		return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
	}

	return r, nil
}

func queryInts(c echo.Context, key string) ([]int64, error) {
	ids := make([]int64, 0)
	idsStr := c.QueryParams()[key]
	for _, v := range idsStr {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, echo.ErrBadRequest.WithInternal(err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func queryIntOptional(c echo.Context, key string) (int, error) {
	str := c.QueryParam(key)
	if str == "" {
		return 0, nil
	}

	number, err := strconv.Atoi(str)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func queryBoolOptional(c echo.Context, key string) (bool, error) {
	str := c.QueryParam(key)
	if str == "" {
		return false, nil
	}

	bool, err := strconv.ParseBool(str)
	if err != nil {
		return false, echo.ErrBadRequest.WithInternal(err)
	}

	return bool, nil
}
