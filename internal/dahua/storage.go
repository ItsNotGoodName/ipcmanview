package dahua

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
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

func normalizeStorageDestination(arg models.DahuaStorageDestination, create bool) models.DahuaStorageDestination {
	arg.Name = strings.TrimSpace(arg.Name)
	arg.ServerAddress = strings.TrimSpace(arg.ServerAddress)

	if arg.Port == 0 {
		switch arg.Storage {
		case models.StorageFTP:
			arg.Port = 21
		case models.StorageSFTP:
			arg.Port = 22
		}
	}

	if create {
		if arg.Name == "" {
			arg.Name = arg.ServerAddress + ":" + strconv.FormatInt(arg.Port, 10)
		}
	}

	return arg
}

func CreateStorageDestination(ctx context.Context, db repo.DB, arg models.DahuaStorageDestination) (int64, error) {
	arg = normalizeStorageDestination(arg, true)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return 0, err
	}

	return db.CreateDahuaStorageDestination(ctx, repo.CreateDahuaStorageDestinationParams{
		Name:            arg.Name,
		Storage:         arg.Storage,
		ServerAddress:   arg.ServerAddress,
		Port:            arg.Port,
		Username:        arg.Username,
		Password:        arg.Password,
		RemoteDirectory: arg.RemoteDirectory,
	})
}

func UpdateStorageDestination(ctx context.Context, db repo.DB, arg models.DahuaStorageDestination) (models.DahuaStorageDestination, error) {
	arg = normalizeStorageDestination(arg, false)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return models.DahuaStorageDestination{}, err
	}

	res, err := db.UpdateDahuaStorageDestination(ctx, repo.UpdateDahuaStorageDestinationParams{
		Name:            arg.Name,
		Storage:         arg.Storage,
		ServerAddress:   arg.ServerAddress,
		Port:            arg.Port,
		Username:        arg.Username,
		Password:        arg.Password,
		RemoteDirectory: arg.RemoteDirectory,
		ID:              arg.ID,
	})
	if err != nil {
		return models.DahuaStorageDestination{}, err
	}

	return res.Convert(), nil
}

func DeleteStorageDestination(ctx context.Context, db repo.DB, id int64) error {
	return db.DeleteDahuaStorageDestination(ctx, id)
}
