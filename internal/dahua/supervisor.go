package dahua

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
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

	bus.OnDahuaCameraCreated(func(ctx context.Context, evt event.DahuaCameraCreated) error {
		_, err := s.GetOrCreateWorker(ctx, evt.CameraID)
		if err != nil {
			return err
		}

		return nil
	})

	bus.OnDahuaCameraUpdated(func(ctx context.Context, evt event.DahuaCameraUpdated) error {
		worker, err := s.GetOrCreateWorker(ctx, evt.CameraID)
		if err != nil {
			return err
		}

		return worker.Restart(ctx)
	})

	bus.OnDahuaCameraDeleted(func(ctx context.Context, evt event.DahuaCameraDeleted) error {
		s.DeleteWorker(evt.CameraID)
		return nil
	})
}

// Sync pulls cameras from database and merges them with the current supervised camers.
// New cameras are created, existing cameras are restarted, and the rest are deleted.
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

func (s *Supervisor) ClientCGI(ctx context.Context, cameraID int64) (dahuacgi.Client, error) {
	worker, err := s.GetOrCreateWorker(ctx, cameraID)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-worker.doneC:
		return nil, ErrWorkerClosed
	case cam, ok := <-worker.camC:
		if !ok {
			return nil, ErrWorkerClosed
		}
		// TODO: reuse connection
		return dahuacgi.NewConn(http.Client{}, cam.Address, cam.Username, cam.Password), nil
	}
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
	camC        chan Camera

	rpcC chan<- workerRPCRequest

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
	camC := make(chan Camera)
	camRefetchC := make(chan struct{}, 1)
	rpcRestartC := make(chan struct{}, 1)
	eventRestartC := make(chan struct{}, 1)
	super.Add(NewWorkerCamera(cameraID, db, camC, camRefetchC, []chan<- struct{}{rpcRestartC, eventRestartC}))

	// RPC worker
	rpcC := make(chan workerRPCRequest)
	super.Add(NewWorkerRPC(camC, rpcC, rpcRestartC))

	// Event worker
	super.Add(NewWorkerEvent(cameraID, db, camC, eventRestartC))

	return SupervisorWorker{
		Supervisor:  super,
		CameraID:    cameraID,
		CamRefetchC: camRefetchC,
		camC:        camC,
		rpcC:        rpcC,
		doneC:       doneC,
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
	case w.CamRefetchC <- struct{}{}:
	default:
	}

	select {
	case <-w.doneC:
		return ErrWorkerClosed
	default:
		return nil
	}
}
