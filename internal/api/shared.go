package api

import (
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	echo "github.com/labstack/echo/v4"
)

func ParamID(c echo.Context) (int64, error) {
	number, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, echo.ErrBadRequest.WithInternal(err)
	}

	return number, nil
}

func QueryTimeRange(start, end string) (models.TimeRange, error) {
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
