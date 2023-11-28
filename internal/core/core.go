package core

import (
	"errors"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func NewLocation(location string) (models.Location, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return models.Location{}, err
	}

	return models.Location{
		Location: loc,
	}, nil
}

func NewTimeRange(start, end time.Time) (models.TimeRange, error) {
	if end.Before(start) {
		return models.TimeRange{}, errors.New("end time before start time")
	}

	return models.TimeRange{
		Start: start,
		End:   end,
	}, nil
}
