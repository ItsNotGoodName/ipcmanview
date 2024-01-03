package dahua

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateCredential(ctx context.Context, db repo.DB, arg repo.CreateDahuaCredentialParams) (int64, error) {
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

	if arg.Name == "" {
		arg.Name = arg.ServerAddress + ":" + strconv.FormatInt(arg.Port, 10)
	}

	if arg.ServerAddress == "" {
		return 0, fmt.Errorf("server address cannot be empty")
	}

	return db.CreateDahuaCredential(ctx, arg)
}

func UpdateCredential(ctx context.Context, db repo.DB, cred repo.DahuaCredential, arg repo.UpdateDahuaCredentialParams) (repo.DahuaCredential, error) {
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

	if arg.Name == "" {
		arg.Name = arg.ServerAddress + ":" + strconv.FormatInt(arg.Port, 10)
	}

	if arg.ServerAddress == "" {
		return repo.DahuaCredential{}, fmt.Errorf("server address cannot be empty")
	}

	return db.UpdateDahuaCredential(ctx, arg)
}

func DeleteCredential(ctx context.Context, db repo.DB, id int64) error {
	return db.DeleteDahuaCredential(ctx, id)
}
