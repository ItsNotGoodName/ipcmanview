package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

var ErrSubscriptionClosed = errors.New("subscription closed")

type PubSub interface {
	SubscribeDahuaEvents(ctx context.Context, cameraIDs []int64) (<-chan models.EventDahuaCameraEvent, error)
}

// ---------- Stream

func useStream(c echo.Context) *json.Encoder {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response())
}

func sendStreamError(c echo.Context, enc *json.Encoder, err error) error {
	str := err.Error()
	if encodeErr := enc.Encode(models.StreamPayload{
		OK:      false,
		Message: &str,
	}); encodeErr != nil {
		return errors.Join(encodeErr, err)
	}

	c.Response().Flush()

	return err
}

func sendStream(c echo.Context, enc *json.Encoder, data any) error {
	err := enc.Encode(models.StreamPayload{
		OK:   true,
		Data: data,
	})
	if err != nil {
		return sendStreamError(c, enc, err)
	}

	c.Response().Flush()

	return nil
}

// ---------- Queries

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

func queryDahuaScanRange(startStr, endStr string) (models.TimeRange, error) {
	end := time.Now()
	start := end.Add(-dahua.MaxScanPeriod)
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
		}
	}

	res, err := core.NewTimeRange(start, end)
	if err != nil {
		return models.TimeRange{}, echo.ErrBadRequest.WithInternal(err)
	}

	return res, nil
}

var encoder = schema.NewEncoder()

var decoder = schema.NewDecoder()

func ParseForm(c echo.Context, form any) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	if err := decoder.Decode(form, c.Request().PostForm); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return nil
}

func ParseLocation(location string) (*time.Location, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, echo.ErrBadRequest.WithInternal(err)
	}
	return loc, nil
}

func ParseQuery(c echo.Context, query any) error {
	if err := decoder.Decode(query, c.Request().URL.Query()); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return nil
}

func ValidateStruct(src any) error {
	err := validate.Validate.Struct(src)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	return err
}

func NewQuery(src any) url.Values {
	query := make(url.Values)
	err := encoder.Encode(src, query)
	if err != nil {
		panic(err)
	}
	return query
}

func FormatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}

func PathID(c echo.Context) (int64, error) {
	number, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func QueryInt(c echo.Context, key string) (int64, error) {
	str := c.QueryParam(key)
	if str == "" {
		return 0, echo.ErrBadRequest
	}

	number, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}
