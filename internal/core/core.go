package core

import (
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
