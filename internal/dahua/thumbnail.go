package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/spf13/afero"
)

type ThumbnailForeignKeys struct {
	FileID            int64
	EmailAttachmentID int64
}

type Thumbnail struct {
	AferoFile
	repo.DahuaThumbnail
}

func CreateThumbnail(ctx context.Context, db repo.DB, afs afero.Fs, fk ThumbnailForeignKeys, width, height int64, aferoFileName string) (Thumbnail, error) {
	thumbnail, err := db.CreateDahuaThumbnail(ctx, repo.CreateDahuaThumbnailParams{
		EmailAttachmentID: core.Int64ToNullInt64(fk.EmailAttachmentID),
		FileID:            core.Int64ToNullInt64(fk.FileID),
		Width:             int64(width),
		Height:            int64(height),
	})
	if err != nil {
		return Thumbnail{}, err
	}

	aferoFile, err := CreateAferoFile(ctx, db, afs, AferoForeignKeys{ThumbnailID: thumbnail.ID}, aferoFileName)
	if err != nil {
		return Thumbnail{}, err
	}

	return Thumbnail{
		AferoFile:      aferoFile,
		DahuaThumbnail: thumbnail,
	}, err
}
