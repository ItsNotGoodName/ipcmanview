package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateDevice(ctx context.Context, db repo.DB, bus *core.Bus, args repo.CreateDahuaDeviceParams) error {
	id, err := db.CreateDahuaDevice(ctx, args, NewFileCursor())
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

func UpdateDevice(ctx context.Context, db repo.DB, bus *core.Bus, args repo.UpdateDahuaDeviceParams) error {
	_, err := db.UpdateDahuaDevice(ctx, args)
	if err != nil {
		return err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, args.ID)
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
