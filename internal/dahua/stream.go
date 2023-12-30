package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func UpdateStream(ctx context.Context, db repo.DB, stream repo.DahuaStream, arg repo.UpdateDahuaStreamParams) (repo.DahuaStream, error) {
	return db.UpdateDahuaStream(ctx, arg)
}
