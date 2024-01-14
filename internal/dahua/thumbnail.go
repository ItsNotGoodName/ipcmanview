package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type ThumbnailForeignKeys struct {
	FileID            int64
	EmailAttachmentID int64
}

func CreateThumbnail(ctx context.Context, db repo.DB, fk ThumbnailForeignKeys, width, height int64) (repo.DahuaThumbnail, error) {
	return db.CreateDahuaThumbnail(ctx, repo.CreateDahuaThumbnailParams{
		EmailAttachmentID: core.Int64ToNullInt64(fk.EmailAttachmentID),
		FileID:            core.Int64ToNullInt64(fk.FileID),
		Width:             int64(width),
		Height:            int64(height),
	})
}
