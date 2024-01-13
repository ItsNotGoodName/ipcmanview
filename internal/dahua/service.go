package dahua

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

var ErrFileServiceConflict = fmt.Errorf("file service conflict")

func NewAferoService(db repo.DB, afs afero.Fs) AferoService {
	return AferoService{
		interval: 8 * time.Hour,
		db:       db,
		afs:      afs,
		queueC:   make(chan struct{}, 1),
	}
}

// AferoService handles deleting orphan afero files.
type AferoService struct {
	interval time.Duration
	db       repo.DB
	afs      afero.Fs
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
	_, err := DeleteOrphanAferoFiles(ctx, s.db, s.afs)
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

type fileServiceReq struct {
	ids  []int64
	resC chan<- fileServiceRes
}

type fileServiceRes struct {
	downloaded int
	err        error
}

type fileServiceFilterReq struct {
	filter repo.DahuaFileFilter
	resC   chan<- fileServiceFilterRes
}

type fileServiceFilterRes struct {
	downloaded int
	err        error
}

func NewFileService(db repo.DB, afs afero.Fs, store *Store) FileService {
	return FileService{
		db:        db,
		afs:       afs,
		store:     store,
		req:       make(chan fileServiceReq),
		filterReq: make(chan fileServiceFilterReq),
	}
}

// FileService handles downloading files.
type FileService struct {
	db        repo.DB
	afs       afero.Fs
	store     *Store
	req       chan fileServiceReq
	filterReq chan fileServiceFilterReq
}

func (s FileService) String() string {
	return "dahua.FileService"
}

func (s FileService) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, s.serve(ctx))
}

func (s FileService) serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-s.filterReq:
			downloaded, err := FileLocalDownloadByFilter(ctx, s.db, s.afs, s.store, req.filter)
			req.resC <- fileServiceFilterRes{
				downloaded: downloaded,
				err:        err,
			}
		case req := <-s.req:
			downloaded := 0
			var firstErr error
			for _, id := range req.ids {
				err := s.download(ctx, id)
				if err != nil {
					if firstErr == nil {
						firstErr = err
					}
					log.Err(err).Msg("Failed to download file")
				} else {
					downloaded++
				}
			}

			req.resC <- fileServiceRes{
				downloaded: downloaded,
				err:        firstErr,
			}
		}
	}
}

func (s FileService) download(ctx context.Context, id int64) error {
	file, err := s.db.GetDahuaFile(ctx, id)
	if err != nil {
		return err
	}

	device, err := s.db.GetDahuaDevice(ctx, file.DeviceID)
	if err != nil {
		return err
	}
	client := s.store.Client(ctx, device.Convert().DahuaConn)

	return FileLocalDownload(ctx, s.db, s.afs, client, file.ID, file.FilePath, file.Type)
}

func (s FileService) DownloadByFilter(ctx context.Context, filter repo.DahuaFileFilter) (int, error) {
	resC := make(chan fileServiceFilterRes, 1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case s.filterReq <- fileServiceFilterReq{
		filter: filter,
		resC:   resC,
	}:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case res := <-resC:
			return res.downloaded, res.err
		}
	default:
		return 0, ErrFileServiceConflict
	}
}

func (s FileService) Download(ctx context.Context, fileIDs ...int64) (int, error) {
	resC := make(chan fileServiceRes, 1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case s.req <- fileServiceReq{
		ids:  fileIDs,
		resC: resC,
	}:
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case res := <-resC:
			return res.downloaded, res.err
		}
	default:
		return 0, ErrFileServiceConflict
	}
}
