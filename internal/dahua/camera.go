package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func DeleteCamera(ctx context.Context, id int64, db repo.DB, dahuaBus *dahuacore.Bus) error {
	if err := db.DeleteDahuaCamera(ctx, id); err != nil {
		return err
	}
	dahuaBus.CameraDeleted(id)
	return nil
}
