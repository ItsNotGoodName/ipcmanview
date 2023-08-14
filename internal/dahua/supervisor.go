package dahua

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/thejerf/suture/v4"
)

type supervisorWorkers []SupervisorWorker

func (s supervisorWorkers) Index(cameraID int64) int {
	for i := range s {
		if s[i].CameraID != cameraID {
			return i
		}
	}

	return -1
}

type Supervisor struct {
	*suture.Supervisor
	db qes.Querier

	mu      sync.Mutex
	workers supervisorWorkers
	tokens  []suture.ServiceToken
}

func NewSupervisor(db qes.Querier) *Supervisor {
	return &Supervisor{
		db:         db,
		Supervisor: suture.NewSimple("dahua.Supervisor"),
	}
}

func (s *Supervisor) Register(bus event.Bus) {
	bus.OnBacklog(func(ctx context.Context) error {
		return s.Sync(ctx)
	})

	bus.OnDahuaCameraDeleted(func(ctx context.Context, evt event.DahuaCameraDeleted) error {
		s.DeleteWorker(evt.CameraID)
		return nil
	})

	bus.OnDahuaCameraUpdated(func(ctx context.Context, evt event.DahuaCameraUpdated) error {
		worker, err := s.GetOrCreateWorker(ctx, evt.CameraID)
		if err != nil {
			return err
		}

		return worker.Restart(ctx)
	})
}

func (s *Supervisor) Sync(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cameras, err := DB.CameraList(ctx, s.db)
	if err != nil {
		return err
	}

	old := make(map[int64]int)
	for i := range s.workers {
		old[s.workers[i].CameraID] = i
	}

	for _, cam := range cameras {
		i, found := old[cam.ID]
		if found {
			// Refresh
			if err := s.workers[i].Restart(ctx); err != nil {
				return err
			}
		} else {
			// Create
			s.createWorker(cam.ID)
		}
	}

	return nil
}

func (s *Supervisor) ClientRPC(ctx context.Context, cameraID int64) (dahuarpc.Client, error) {
	return s.GetOrCreateWorker(ctx, cameraID)
}

func (s *Supervisor) GetOrCreateWorker(ctx context.Context, cameraID int64) (SupervisorWorker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.workers {
		if s.workers[i].CameraID != cameraID {
			continue
		}

		// Found
		return s.workers[i], nil
	}

	if err := DB.CameraExists(ctx, s.db, cameraID); err != nil {
		return SupervisorWorker{}, err
	}

	// Create

	return s.createWorker(cameraID), nil
}

func (s *Supervisor) createWorker(cameraID int64) SupervisorWorker {
	worker := newSupervisorWorker(cameraID, s.db)
	token := s.Add(worker)

	s.workers = append(s.workers, worker)
	s.tokens = append(s.tokens, token)

	return worker
}

func (s *Supervisor) DeleteWorker(cameraID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.workers.Index(cameraID)
	if index == -1 {
		return
	}

	s.deleteWorker(cameraID)
}

func (s *Supervisor) deleteWorker(cameraID int64) {
	workers := []SupervisorWorker{}
	tokens := []suture.ServiceToken{}
	for i := range s.workers {
		if s.workers[i].CameraID != cameraID {
			workers = append(workers, s.workers[i])
			tokens = append(tokens, s.tokens[i])
			continue
		}

		s.Supervisor.Remove(s.tokens[i])
	}

	s.workers = workers
	s.tokens = tokens
}

type SupervisorWorker struct {
	*suture.Supervisor
	CameraID int64

	CamRefetchC chan struct{}
	camC        chan models.DahuaCamera

	rpcRestartC chan struct{}
	rpcC        chan<- workerRPCRequest

	eventRestartC chan struct{}

	doneC chan struct{}
}

func newSupervisorWorker(cameraID int64, db qes.Querier) SupervisorWorker {
	super := suture.New(fmt.Sprintf("dahua.SupervisorWorker@camera-%d", cameraID), suture.Spec{
		DontPropagateTermination: true,
	})

	// Done worker
	doneC := make(chan struct{})
	super.Add(NewWorkerDone(cameraID, doneC))

	// Camera worker
	camC := make(chan models.DahuaCamera)
	camRefetchC := make(chan struct{}, 1)
	super.Add(NewWorkerCamera(cameraID, db, camC, camRefetchC))

	// RPC worker
	rpcC := make(chan workerRPCRequest)
	rpcRestartC := make(chan struct{}, 1)
	super.Add(NewWorkerRPC(camC, rpcC, rpcRestartC))

	// Event worker
	eventRestartC := make(chan struct{}, 1)
	super.Add(NewWorkerEvent(cameraID, db, camC, eventRestartC))

	return SupervisorWorker{
		Supervisor:    super,
		CameraID:      cameraID,
		CamRefetchC:   camRefetchC,
		camC:          camC,
		rpcRestartC:   rpcRestartC,
		rpcC:          rpcC,
		eventRestartC: eventRestartC,
		doneC:         doneC,
	}
}

func (w SupervisorWorker) RPC(ctx context.Context) (dahuarpc.RequestBuilder, error) {
	resC := make(chan workerRPCResponse)
	select {
	case <-ctx.Done():
		return dahuarpc.RequestBuilder{}, ctx.Err()
	case <-w.doneC:
		return dahuarpc.RequestBuilder{}, ErrWorkerClosed
	case w.rpcC <- workerRPCRequest{ctx, resC}:
		res := <-resC
		if res.err != nil {
			return dahuarpc.RequestBuilder{}, res.err
		}

		return res.rpc, nil
	}
}

func (w SupervisorWorker) Restart(ctx context.Context) error {
	// Refetch camera from database
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.doneC:
		return ErrWorkerClosed
	case w.CamRefetchC <- struct{}{}:
	}

	// Restart rpc client
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.doneC:
		return ErrWorkerClosed
	case w.rpcRestartC <- struct{}{}:
	}

	// restart event client
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.doneC:
		return ErrWorkerClosed
	case w.eventRestartC <- struct{}{}:
	}

	return nil
}
