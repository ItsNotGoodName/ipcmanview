package dahua

import (
	"context"
	"fmt"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/spf13/afero"
)

// SyncAferoFile deletes the file from the database if it does not exist in the file system.
func SyncAferoFile(ctx context.Context, db repo.DB, fs afero.Fs, aferoFile repo.DahuaAferoFile, err error) (bool, error) {
	if err != nil {
		if repo.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return syncAferoFile(ctx, db, fs, aferoFile)
}

func syncAferoFile(ctx context.Context, db repo.DB, fs afero.Fs, aferoFile repo.DahuaAferoFile) (bool, error) {
	_, err := fs.Stat(aferoFile.Name)
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

// DeleteAllOrphanAferoFile deletes unreferenced afero files.
func DeleteAllOrphanAferoFile(ctx context.Context, db repo.DB, fs afero.Fs) (int, error) {
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
			err := fs.Remove(f.Name)
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
