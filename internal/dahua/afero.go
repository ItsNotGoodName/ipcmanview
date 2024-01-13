package dahua

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/spf13/afero"
)

const AferoEchoRoute = "/v1/dahua-afero-files/*"
const AferoEchoRoutePrefix = "/v1/dahua-afero-files"

func AferoFileURI(name string) string {
	return "/v1/dahua-afero-files/" + name
}

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

// DeleteOrphanAferoFiles deletes unreferenced afero files.
func DeleteOrphanAferoFiles(ctx context.Context, db repo.DB, fs afero.Fs) (int, error) {
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

func NewAferoService(db repo.DB, fs afero.Fs) AferoService {
	return AferoService{
		interval: 8 * time.Hour,
		db:       db,
		fs:       fs,
		queueC:   make(chan struct{}, 1),
	}
}

// AferoService handles deleting orphan afero files.
type AferoService struct {
	interval time.Duration
	db       repo.DB
	fs       afero.Fs
	queueC   chan struct{}
}

func (s AferoService) String() string {
	return "dahua.AferoService"
}

func (s AferoService) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, s.serve(ctx))
}

func (s AferoService) serve(ctx context.Context) error {
	t := time.NewTicker(s.interval)
	defer t.Stop()

	if err := s.run(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.queueC:
			if err := s.run(ctx); err != nil {
				return err
			}
		case <-t.C:
			if err := s.run(ctx); err != nil {
				return err
			}
		}
	}
}

func (s AferoService) run(ctx context.Context) error {
	_, err := DeleteOrphanAferoFiles(ctx, s.db, s.fs)
	if err != nil {
		return err
	}

	return nil
}

func (s AferoService) Queue() {
	select {
	case s.queueC <- struct{}{}:
	default:
	}
}
