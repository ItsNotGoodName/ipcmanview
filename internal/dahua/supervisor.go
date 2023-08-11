package dahua

import (
	"context"
	"errors"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/thejerf/suture/v4"
)

type Supervisor struct {
	*suture.Supervisor
	db       qes.Querier
	mu       sync.Mutex
	registry []SupervisorWorker
}

type SupervisorWorker struct {
	doneC  chan<- struct{}
	worker Worker
	token  suture.ServiceToken
}

func NewSupervisor(db qes.Querier, bus *event.Bus) *Supervisor {
	return &Supervisor{
		db:         db,
		Supervisor: suture.NewSimple("dahua.Supervisor"),
	}
}

func (s *Supervisor) GetGenRPC(ctx context.Context, cameraID int64) (dahua.GenRPC, error) {
	return s.GetWorker(ctx, cameraID)
}

func (s *Supervisor) GetWorker(ctx context.Context, cameraID int64) (Worker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.registry {
		if s.registry[i].worker.CameraID != cameraID {
			continue
		}

		// Found
		return s.registry[i].worker, nil
	}

	if err := DB.CameraExists(ctx, s.db, cameraID); err != nil {
		return Worker{}, err
	}

	// Create
	doneC := make(chan struct{})
	worker := NewWorker(s.db, cameraID, doneC)
	workerToken := s.Add(newSupervisorWorker(s, worker))

	s.registry = append(s.registry, SupervisorWorker{
		doneC:  doneC,
		worker: worker,
		token:  workerToken,
	})

	return worker, nil
}

func (s *Supervisor) DeleteWorker(cameraID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.registry {
		if s.registry[i].worker.CameraID == cameraID {
			s.deleteWorker(cameraID)
			return
		}
	}
}

// deleteWorker deletes worker by cameraID.
func (s *Supervisor) deleteWorker(cameraID int64) {
	registry := []SupervisorWorker{}
	for i := range s.registry {
		if s.registry[i].worker.CameraID != cameraID {
			registry = append(registry, s.registry[i])
			continue
		}

		// Delete
		s.Supervisor.Remove(s.registry[i].token)
		close(s.registry[i].doneC)
	}

	s.registry = registry
}

// supervisorWorker deletes the worker from supervisor when shutting down.
type supervisorWorker struct {
	worker     Worker
	supervisor *Supervisor
}

func newSupervisorWorker(supervisor *Supervisor, worker Worker) *supervisorWorker {
	return &supervisorWorker{
		worker:     worker,
		supervisor: supervisor,
	}
}

func (s supervisorWorker) Serve(ctx context.Context) error {
	err := s.worker.Serve(ctx)

	if errors.Is(err, suture.ErrDoNotRestart) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		s.supervisor.DeleteWorker(s.worker.CameraID)
	}

	return err
}
