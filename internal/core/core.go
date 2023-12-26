package core

import (
	"errors"
	"strings"
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

func StorageFromFilePath(filePath string) models.Storage {
	if strings.HasPrefix(filePath, "sftp://") {
		return models.StorageSFTP
	}
	if strings.HasPrefix(filePath, "ftp://") {
		return models.StorageFTP
	}
	// if strings.HasPrefix(filePath, "nfs://") {
	// 	return models.StorageNFS
	// }
	// if strings.HasPrefix(filePath, "smb://") {
	// 	return models.StorageSMB
	// }
	return models.StorageLocal
}
