package dahua

import (
	"context"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateDevice(ctx context.Context, db repo.DB, bus *core.Bus, arg repo.CreateDahuaDeviceParams) error {
	arg.Address.URL = toHTTPURL(arg.Address.URL)
	arg.Name = strings.TrimSpace(arg.Name)

	id, err := db.CreateDahuaDevice(ctx, arg, NewFileCursor())
	if err != nil {
		return err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, id)
	if err != nil {
		return err
	}
	bus.EventDahuaDeviceCreated(models.EventDahuaDeviceCreated{
		Device: dbDevice.Convert(),
	})

	return nil
}

func UpdateDevice(ctx context.Context, db repo.DB, bus *core.Bus, device models.DahuaDevice, arg repo.UpdateDahuaDeviceParams) error {
	arg.Address.URL = toHTTPURL(arg.Address.URL)
	arg.Name = strings.TrimSpace(arg.Name)

	_, err := db.UpdateDahuaDevice(ctx, arg)
	if err != nil {
		return err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, arg.ID)
	if err != nil {
		return err
	}
	bus.EventDahuaDeviceUpdated(models.EventDahuaDeviceUpdated{
		Device: dbDevice.Convert(),
	})

	return nil
}

func DeleteDevice(ctx context.Context, db repo.DB, bus *core.Bus, id int64) error {
	if err := db.DeleteDahuaDevice(ctx, id); err != nil {
		return err
	}
	bus.EventDahuaDeviceDeleted(models.EventDahuaDeviceDeleted{
		DeviceID: id,
	})
	return nil
}
