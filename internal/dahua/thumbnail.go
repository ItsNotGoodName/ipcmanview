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

type Thumbnail struct {
	aferoFile
	repo.DahuaThumbnail
}

func CreateThumbnail(ctx context.Context, fk ThumbnailForeignKeys, width, height int64, aferoFileName string) (Thumbnail, error) {
	thumbnail, err := app.DB.C().DahuaCreateThumbnail(ctx, repo.DahuaCreateThumbnailParams{
		EmailAttachmentID: core.Int64ToNullInt64(fk.EmailAttachmentID),
		FileID:            core.Int64ToNullInt64(fk.FileID),
		Width:             int64(width),
		Height:            int64(height),
	})
	if err != nil {
		return Thumbnail{}, err
	}

	aferoFile, err := createAferoFile(ctx, aferoForeignKeys{ThumbnailID: thumbnail.ID}, aferoFileName)
	if err != nil {
		return Thumbnail{}, err
	}

	return Thumbnail{
		aferoFile:      aferoFile,
		DahuaThumbnail: thumbnail,
	}, err
}
