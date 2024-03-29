package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/labstack/echo/v4"
)

func paramID(c echo.Context) (int64, error) {
	number, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}
	return number, nil
}

func assertDahuaLevel(c echo.Context, s *Server, deviceID int64, level models.DahuaPermissionLevel) error {
	ok, err := dahua.Level(c.Request().Context(), deviceID, level)
	if err != nil {
		if core.IsNotFound(err) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}
	if !ok {
		return echo.ErrForbidden
	}
	return nil
}

func useDahuaClient(c echo.Context, s *Server, deviceID int64) (dahua.Client, error) {
	client, err := dahua.GetClient(c.Request().Context(), deviceID)
	if err != nil {
		if core.IsNotFound(err) {
			return dahua.Client{}, echo.ErrNotFound.WithInternal(err)
		}
		return dahua.Client{}, err
	}
	return client, nil
}

// ---------- Stream

func newStream(c echo.Context) *json.Encoder {
	c.Response().Header().Set(echo.HeaderContentType, "application/x-ndjson")
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
