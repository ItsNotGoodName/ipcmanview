package dahua

import (
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

var Storage = []models.Storage{
	models.StorageFTP,
	models.StorageSFTP,
}

func ParseStorage(storage string) (models.Storage, error) {
	switch storage {
	case string(models.StorageFTP):
		return models.StorageFTP, nil
	case string(models.StorageSFTP):
		return models.StorageSFTP, nil
	}
	return "", fmt.Errorf("storage not supported: %s", storage)
}
