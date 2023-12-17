package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func DeleteCamera(ctx context.Context, db repo.DB, bus *core.Bus, id int64) error {
	if err := db.DeleteDahuaCamera(ctx, id); err != nil {
		return err
	}
	bus.EventDahuaCameraDeleted(models.EventDahuaCameraDeleted{
		CameraID: id,
	})
	return nil
}
