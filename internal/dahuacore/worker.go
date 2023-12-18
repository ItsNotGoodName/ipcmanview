package dahuacore

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/thejerf/suture/v4"
)

func DefaultWorkerBuilder(hooks EventHooks, bus *core.Bus, store *Store) WorkerBuilder {
	return func(ctx context.Context, super *suture.Supervisor, camera models.DahuaConn) ([]suture.ServiceToken, error) {
		var tokens []suture.ServiceToken

		{
			worker := NewEventWorker(camera, hooks)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		{
			worker := NewCoaxialWorker(bus, camera.ID, store.Conn(ctx, camera).RPC)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		return tokens, nil
	}
}

type WorkerBuilder = func(ctx context.Context, super *suture.Supervisor, camera models.DahuaConn) ([]suture.ServiceToken, error)

type WorkerStore struct {
	super *suture.Supervisor
	build WorkerBuilder

	workersMu sync.Mutex
	workers   map[int64]workerData
}

type workerData struct {
	camera models.DahuaConn
	tokens []suture.ServiceToken
}

func NewWorkerStore(super *suture.Supervisor, build WorkerBuilder) *WorkerStore {
	return &WorkerStore{
		super:     super,
		build:     build,
		workersMu: sync.Mutex{},
		workers:   make(map[int64]workerData),
	}
}

func (s *WorkerStore) create(ctx context.Context, camera models.DahuaConn) error {
	tokens, err := s.build(ctx, s.super, camera)
	if err != nil {
		return err
	}

	s.workers[camera.ID] = workerData{
		camera: camera,
		tokens: tokens,
	}

	return nil
}

func (s *WorkerStore) Create(ctx context.Context, camera models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	_, found := s.workers[camera.ID]
	if found {
		return fmt.Errorf("workers already exists for camera by ID: %d", camera.ID)
	}

	return s.create(ctx, camera)
}

func (s *WorkerStore) Update(ctx context.Context, camera models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	worker, found := s.workers[camera.ID]
	if !found {
		return fmt.Errorf("workers not found for camera by ID: %d", camera.ID)
	}
	if ConnEqual(worker.camera, camera) {
		return nil
	}

	for _, st := range worker.tokens {
		s.super.Remove(st)
	}

	return s.create(ctx, camera)
}

func (s *WorkerStore) Delete(id int64) error {
	s.workersMu.Lock()
	worker, found := s.workers[id]
	if !found {
		s.workersMu.Unlock()
		return fmt.Errorf("workers not found for camera by ID: %d", id)
	}

	for _, token := range worker.tokens {
		s.super.Remove(token)
	}
	delete(s.workers, id)
	s.workersMu.Unlock()
	return nil
}

func (e *WorkerStore) Register(bus *core.Bus) {
	bus.OnEventDahuaCameraCreated(func(ctx context.Context, evt models.EventDahuaCameraCreated) error {
		return e.Create(ctx, evt.Camera.DahuaConn)
	})
	bus.OnEventDahuaCameraUpdated(func(ctx context.Context, evt models.EventDahuaCameraUpdated) error {
		return e.Update(ctx, evt.Camera.DahuaConn)
	})
	bus.OnEventDahuaCameraDeleted(func(ctx context.Context, evt models.EventDahuaCameraDeleted) error {
		return e.Delete(evt.CameraID)
	})
}

type CameraStore interface {
	ListConn(ctx context.Context) ([]models.DahuaConn, error)
}

func (w *WorkerStore) Bootstrap(ctx context.Context, cameraStore CameraStore, store *Store) error {
	cameras, err := cameraStore.ListConn(ctx)
	if err != nil {
		return err
	}
	conns := store.ConnList(ctx, cameras)
	for _, conn := range conns {
		if err := w.Create(ctx, conn.Camera); err != nil {
			return err
		}
	}
	return err
}
