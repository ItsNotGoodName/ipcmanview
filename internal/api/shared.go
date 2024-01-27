package api

import (
	"net/url"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"github.com/gorilla/schema"
	echo "github.com/labstack/echo/v4"
)

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

func DecodeQuery(c echo.Context, dst any) error {
	if err := decoder.Decode(dst, c.Request().URL.Query()); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	return nil
}

func EncodeQuery(src any) url.Values {
	query := make(url.Values)
	err := encoder.Encode(src, query)
	if err != nil {
		panic(err)
	}
	return query
}

func ValidateStruct(src any) error {
	err := validate.Validate.Struct(src)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	return err
}

func ParamID(c echo.Context) (int64, error) {
	number, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func UseTimeRange(start, end string) (models.TimeRange, error) {
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
