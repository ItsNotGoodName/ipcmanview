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

const AferoEchoRoute = "/v1/dahua-afero-files/*"
const AferoEchoRoutePrefix = "/v1/dahua-afero-files"

func AferoFileURI(name string) string {
	return "/v1/dahua-afero-files/" + name
}

func NewAferoFileName(extension string) string {
	uuid := uuid.NewString()
	if extension == "" {
		return uuid
	}
	if strings.HasPrefix(".", extension) {
		return uuid + extension
	}
	return uuid + "." + extension
}

// SyncAferoFile deletes the file from the database if it does not exist in the file system.
func SyncAferoFile(ctx context.Context, db repo.DB, afs afero.Fs, aferoFile repo.DahuaAferoFile, err error) (bool, error) {
	if err != nil {
		if repo.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return syncAferoFile(ctx, db, afs, aferoFile)
}

func syncAferoFile(ctx context.Context, db repo.DB, afs afero.Fs, aferoFile repo.DahuaAferoFile) (bool, error) {
	_, err := afs.Stat(aferoFile.Name)
	if err == nil {
		return true, nil
	}
	if !os.IsNotExist(err) {
		return false, err
	}

	err = db.DeleteDahuaAferoFile(ctx, aferoFile.ID)
	if err != nil {
		return false, err
	}

	return false, nil
}

// DeleteOrphanAferoFiles deletes unreferenced afero files.
func DeleteOrphanAferoFiles(ctx context.Context, db repo.DB, afs afero.Fs) (int, error) {
	deleted := 0

	var first repo.DahuaAferoFile
	for {
		files, err := db.OrphanListDahuaAferoFile(ctx, 20)
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

			err = db.DeleteDahuaAferoFile(ctx, f.ID)
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
	FileThumbnailID   int64
	EmailAttachmentID int64
}

func CreateAferoFile(ctx context.Context, db repo.DB, afs afero.Fs, fileName string, key AferoForeignKeys) (AferoFile, error) {
	id, err := db.CreateDahuaAferoFile(ctx, repo.CreateDahuaAferoFileParams{
		FileID:            core.Int64ToNullInt64(key.FileID),
		FileThumbnailID:   core.Int64ToNullInt64(key.FileThumbnailID),
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

func ReadyAferoFile(ctx context.Context, db repo.DB, id int64, file afero.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = db.ReadyDahuaAferoFile(ctx, repo.ReadyDahuaAferoFileParams{
		Size:      stat.Size(),
		CreatedAt: types.NewTime(stat.ModTime()),
		ID:        id,
	})
	if err != nil {
		return err
	}

	return nil
}
