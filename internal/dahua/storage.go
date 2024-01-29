package dahua

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
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
			arg.Name = core.Address(arg.ServerAddress, int(arg.Port))
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

	return db.DahuaCreateStorageDestination(ctx, repo.DahuaCreateStorageDestinationParams{
		Name:            arg.Name,
		Storage:         arg.Storage,
		ServerAddress:   arg.ServerAddress,
		Port:            arg.Port,
		Username:        arg.Username,
		Password:        arg.Password,
		RemoteDirectory: arg.RemoteDirectory,
	})
}

func UpdateStorageDestination(ctx context.Context, db repo.DB, arg models.DahuaStorageDestination) error {
	arg = normalizeStorageDestination(arg, false)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return err
	}

	_, err = db.DahuaUpdateStorageDestination(ctx, repo.DahuaUpdateStorageDestinationParams{
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
		return err
	}

	return nil
}

func DeleteStorageDestination(ctx context.Context, db repo.DB, id int64) error {
	return db.DahuaDeleteStorageDestination(ctx, id)
}

func TestStorageDestination(ctx context.Context, arg models.DahuaStorageDestination) error {
	switch arg.Storage {
	case models.StorageFTP:
		c, err := ftp.Dial(core.Address(arg.ServerAddress, int(arg.Port)), ftp.DialWithContext(ctx))
		if err != nil {
			return err
		}

		err = c.Login(arg.Username, arg.Password)
		if err != nil {
			return err
		}

		return c.Quit()
	case models.StorageSFTP:
		conn, err := ssh.Dial("tcp", core.Address(arg.ServerAddress, int(arg.Port)), &ssh.ClientConfig{
			User: arg.Username,
			Auth: []ssh.AuthMethod{ssh.Password(arg.Password)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// TODO: check public key
				return nil
			},
		})
		if err != nil {
			return err
		}

		client, err := sftp.NewClient(conn)
		if err != nil {
			return err
		}

		return client.Close()
	default:
		return fmt.Errorf("invalid storage: %s", arg.Storage)
	}
}
