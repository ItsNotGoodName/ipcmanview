package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/thejerf/suture/v4"
)

type ScanSupervisor struct {
	*suture.Supervisor
	workers []scanSupervisorWorker
}

func NewScanSupervisor(db qes.Querier, store Store, scanners int) *ScanSupervisor {
	super := suture.NewSimple("dahua.ScanSupervisor")

	var workers []scanSupervisorWorker
	for i := 0; i < scanners; i++ {
		worker := newScanSupervisorWorker(db, store)
		super.Add(worker)
		workers = append(workers, worker)
	}

	return &ScanSupervisor{
		Supervisor: super,
		workers:    workers,
	}
}

// Scan signals to a worker to start scan if they are not busy.
func (s *ScanSupervisor) Scan() {
	for i := range s.workers {
		if s.workers[i].Scan() {
			return
		}
	}
}

type scanSupervisorWorker struct {
	db    qes.Querier
	store Store
	scanC chan struct{}
}

func (w *scanSupervisorWorker) String() string {
	return "dahua.scanSupervisorWorker"
}

func newScanSupervisorWorker(db qes.Querier, store Store) scanSupervisorWorker {
	return scanSupervisorWorker{
		db:    db,
		store: store,
		scanC: make(chan struct{}),
	}
}

func (s scanSupervisorWorker) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.scanC:
			for {
				ok, err := DB.ScanQueueTaskNext(ctx, s.db, func(ctx context.Context, scanCursorLock ScanCursorLock, queueTask models.DahuaScanQueueTask) error {
					gen, err := s.store.GetGenRPC(ctx, queueTask.CameraID)
					if err != nil {
						return err
					}

					return ScanTaskQueueExecute(ctx, s.db, gen, scanCursorLock, queueTask)
				})
				if err != nil {
					return err
				}
				if !ok {
					break
				}
			}
		}
	}
}

func (s scanSupervisorWorker) Scan() bool {
	select {
	case s.scanC <- struct{}{}:
		return len(s.scanC) == 0
	default:
		return false
	}
}
