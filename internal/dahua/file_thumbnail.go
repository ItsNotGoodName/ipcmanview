package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func CreateFileThumbnail(ctx context.Context, db repo.DB, fileID int64, width, height int64) (repo.DahuaFileThumbnail, error) {
	return db.CreateDahuaFileThumbnail(ctx, repo.CreateDahuaFileThumbnailParams{
		FileID: fileID,
		Width:  int64(width),
		Height: int64(height),
	})
}
