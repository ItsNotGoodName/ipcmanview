package dahua

import (
	"context"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
)

func normalizeDevice(arg models.DahuaDevice, create bool) models.DahuaDevice {
	arg.Name = strings.TrimSpace(arg.Name)
	arg.Address = toHTTPURL(arg.Address)

	if create {
		if arg.Name == "" {
			arg.Name = arg.Address.String()
		}
		if arg.Username == "" {
			arg.Username = "admin"
		}
	}

	return arg
}

func CreateDevice(ctx context.Context, db repo.DB, bus *core.Bus, arg models.DahuaDevice) (models.DahuaDeviceConn, error) {
	arg = normalizeDevice(arg, true)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	now := types.NewTime(time.Now())
	id, err := db.CreateDahuaDevice(ctx, repo.CreateDahuaDeviceParams{
		Name:      arg.Name,
		Address:   types.NewURL(arg.Address),
		Username:  arg.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		CreatedAt: now,
		UpdatedAt: now,
	}, NewFileCursor())
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, id)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}
	device := dbDevice.Convert()

	bus.EventDahuaDeviceCreated(models.EventDahuaDeviceCreated{
		Device: device,
	})

	return device, err
}

func UpdateDevice(ctx context.Context, db repo.DB, bus *core.Bus, arg models.DahuaDevice) (models.DahuaDeviceConn, error) {
	arg = normalizeDevice(arg, false)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	_, err = db.UpdateDahuaDevice(ctx, repo.UpdateDahuaDeviceParams{
		Name:      arg.Name,
		Address:   types.NewURL(arg.Address),
		Username:  arg.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		UpdatedAt: types.NewTime(time.Now()),
		ID:        arg.ID,
	})
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, arg.ID)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}
	device := dbDevice.Convert()

	bus.EventDahuaDeviceUpdated(models.EventDahuaDeviceUpdated{
		Device: device,
	})

	return device, nil
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
