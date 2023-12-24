package core

import (
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func NewLocation(location string) (types.Location, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return types.Location{}, err
	}

	return types.Location{
		Location: loc,
	}, nil
}

func NewTimeRange(start, end time.Time) (models.TimeRange, error) {
	if end.Before(start) {
		return models.TimeRange{}, errors.New("invalid time range: end is before start")
	}

	return models.TimeRange{
		Start: start,
		End:   end,
	}, nil
}
