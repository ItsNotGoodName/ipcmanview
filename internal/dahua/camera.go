package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func DeleteCamera(ctx context.Context, db repo.DB, dahuaBus *dahuacore.Bus, id int64) error {
	if err := db.DeleteDahuaCamera(ctx, id); err != nil {
		return err
	}
	dahuaBus.CameraDeleted(id)
	return nil
}
