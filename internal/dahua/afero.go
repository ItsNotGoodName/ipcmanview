package dahua

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/google/uuid"
	"github.com/spf13/afero"
)

func NewAferoFileName(extension string) string {
	uuid := uuid.NewString()
	if extension == "" {
		return uuid
	}
	if strings.HasPrefix(extension, ".") {
		return uuid + extension
	}
	return uuid + "." + extension
}

// SyncAferoFile deletes the file from the database if it does not exist in the file system.
func SyncAferoFile(ctx context.Context, db sqlite.DB, afs afero.Fs, aferoFileID int64, aferoFileName string) error {
	_, err := afs.Stat(aferoFileName)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	err = db.C().DahuaDeleteAferoFile(ctx, aferoFileID)
	if err != nil {
		return err
	}

	return core.ErrNotFound
}

// DeleteOrphanAferoFiles deletes unreferenced afero files.
func DeleteOrphanAferoFiles(ctx context.Context, db sqlite.DB, afs afero.Fs) (int, error) {
	deleted := 0

	var first repo.DahuaAferoFile
	for {
		files, err := db.C().DahuaOrphanListAferoFiles(ctx, 20)
		if err != nil {
			return deleted, err
		}
		if len(files) == 0 {
			return deleted, nil
		}
		if files[0].ID == first.ID {
			return deleted, fmt.Errorf("repeat afero file: %d", first.ID)
		}
		first = files[0]

		for _, f := range files {
			err := afs.Remove(f.Name)
			if err != nil && !os.IsNotExist(err) {
				return deleted, err
			}

			err = db.C().DahuaDeleteAferoFile(ctx, f.ID)
			if err != nil {
				return deleted, err
			}
			deleted++
		}
	}
}

type AferoFile struct {
	afero.File
	ID   int64
	Name string
}

type AferoForeignKeys struct {
	FileID            int64
	ThumbnailID       int64
	EmailAttachmentID int64
}

// CreateAferoFile creates an afero file in the database and in the file system.
func CreateAferoFile(ctx context.Context, db sqlite.DB, afs afero.Fs, key AferoForeignKeys, fileName string) (AferoFile, error) {
	id, err := db.C().DahuaCreateAferoFile(ctx, repo.DahuaCreateAferoFileParams{
		FileID:            core.Int64ToNullInt64(key.FileID),
		ThumbnailID:       core.Int64ToNullInt64(key.ThumbnailID),
		EmailAttachmentID: core.Int64ToNullInt64(key.EmailAttachmentID),
		Name:              fileName,
		CreatedAt:         types.NewTime(time.Now()),
	})
	if err != nil {
		return AferoFile{}, err
	}

	file, err := afs.Create(fileName)
	if err != nil {
		return AferoFile{}, err
	}

	return AferoFile{
		File: file,
		ID:   id,
		Name: fileName,
	}, nil
}

// ReadyAferoFile sets the afero file to a ready state.
func ReadyAferoFile(ctx context.Context, db sqlite.DB, id int64, file afero.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = db.C().DahuaReadyAferoFile(ctx, repo.DahuaReadyAferoFileParams{
		Size:      stat.Size(),
		CreatedAt: types.NewTime(stat.ModTime()),
		ID:        id,
	})
	if err != nil {
		return err
	}

	return nil
}
