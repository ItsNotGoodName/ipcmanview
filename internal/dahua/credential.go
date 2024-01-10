package dahua

import (
	"context"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
)

func normalizeCredential(arg models.DahuaCredential, create bool) models.DahuaCredential {
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

func CreateCredential(ctx context.Context, db repo.DB, arg models.DahuaCredential) (int64, error) {
	arg = normalizeCredential(arg, true)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return 0, err
	}

	return db.CreateDahuaCredential(ctx, repo.CreateDahuaCredentialParams{
		Name:            arg.Name,
		Storage:         arg.Storage,
		ServerAddress:   arg.ServerAddress,
		Port:            arg.Port,
		Username:        arg.Username,
		Password:        arg.Password,
		RemoteDirectory: arg.RemoteDirectory,
	})
}

func UpdateCredential(ctx context.Context, db repo.DB, arg models.DahuaCredential) (models.DahuaCredential, error) {
	arg = normalizeCredential(arg, false)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return models.DahuaCredential{}, err
	}

	res, err := db.UpdateDahuaCredential(ctx, repo.UpdateDahuaCredentialParams{
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
		return models.DahuaCredential{}, err
	}

	return res.Convert(), nil
}

func DeleteCredential(ctx context.Context, db repo.DB, id int64) error {
	return db.DeleteDahuaCredential(ctx, id)
}
