package dahua

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/google/uuid"
	"github.com/spf13/afero"
)

func newAferoFileName(extension string) string {
	uuid := uuid.NewString()
	if extension == "" {
		return uuid
	}
	if strings.HasPrefix(extension, ".") {
		return uuid + extension
	}
	return uuid + "." + extension
}

// syncAferoFile deletes the file from the database if it does not exist in the file system.
func syncAferoFile(ctx context.Context, aferoFileID int64, aferoFileName string) error {
	_, err := app.AFS.Stat(aferoFileName)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	err = app.DB.C().DahuaDeleteAferoFile(ctx, aferoFileID)
	if err != nil {
		return err
	}

	return core.ErrNotFound
}

// deleteOrphanAferoFiles deletes unreferenced afero files.
func deleteOrphanAferoFiles(ctx context.Context) (int, error) {
	deleted := 0

	var first repo.DahuaAferoFile
	for {
		files, err := app.DB.C().DahuaOrphanListAferoFiles(ctx, 20)
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
			err := app.AFS.Remove(f.Name)
			if err != nil && !os.IsNotExist(err) {
				return deleted, err
			}

			err = app.DB.C().DahuaDeleteAferoFile(ctx, f.ID)
			if err != nil {
				return deleted, err
			}
			deleted++
		}
	}
}

type aferoFile struct {
	afero.File
	ID   int64
	Name string
}

type aferoForeignKeys struct {
	FileID            int64
	ThumbnailID       int64
	EmailAttachmentID int64
}

// createAferoFile creates an afero file in the database and in the file system.
func createAferoFile(ctx context.Context, key aferoForeignKeys, fileName string) (aferoFile, error) {
	id, err := app.DB.C().DahuaCreateAferoFile(ctx, repo.DahuaCreateAferoFileParams{
		FileID:            core.Int64ToNullInt64(key.FileID),
		ThumbnailID:       core.Int64ToNullInt64(key.ThumbnailID),
		EmailAttachmentID: core.Int64ToNullInt64(key.EmailAttachmentID),
		Name:              fileName,
		CreatedAt:         types.NewTime(time.Now()),
	})
	if err != nil {
		return aferoFile{}, err
	}

	file, err := app.AFS.Create(fileName)
	if err != nil {
		return aferoFile{}, err
	}

	return aferoFile{
		File: file,
		ID:   id,
		Name: fileName,
	}, nil
}

// Ready sets the afero file to a ready state.
func (f aferoFile) Ready(ctx context.Context) error {
	stat, err := f.File.Stat()
	if err != nil {
		return err
	}

	_, err = app.DB.C().DahuaReadyAferoFile(ctx, repo.DahuaReadyAferoFileParams{
		Size:      stat.Size(),
		CreatedAt: types.NewTime(stat.ModTime()),
		ID:        f.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
