package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func DeleteDevice(ctx context.Context, db repo.DB, bus *core.Bus, id int64) error {
	if err := db.DeleteDahuaDevice(ctx, id); err != nil {
		return err
	}
	bus.EventDahuaDeviceDeleted(models.EventDahuaDeviceDeleted{
		DeviceID: id,
	})
	return nil
}
