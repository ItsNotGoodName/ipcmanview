package dahua

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/thejerf/suture/v4"
)

type WorkerFactory = func(ctx context.Context, super *suture.Supervisor, device models.Conn) ([]suture.ServiceToken, error)

// WorkerStore manages the lifecycle of workers to devices.
type WorkerStore struct {
	super   *suture.Supervisor
	factory WorkerFactory

	workersMu sync.Mutex
	workers   map[int64]workerData
}

type workerData struct {
	conn   models.Conn
	tokens []suture.ServiceToken
}

func NewWorkerStore(super *suture.Supervisor, factory WorkerFactory) *WorkerStore {
	return &WorkerStore{
		super:     super,
		factory:   factory,
		workersMu: sync.Mutex{},
		workers:   make(map[int64]workerData),
	}
}

func (s *WorkerStore) create(ctx context.Context, conn models.Conn) error {
	tokens, err := s.factory(ctx, s.super, conn)
	if err != nil {
		return err
	}

	s.workers[conn.ID] = workerData{
		conn:   conn,
		tokens: tokens,
	}

	return nil
}

func (s *WorkerStore) Create(ctx context.Context, device models.Conn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	_, found := s.workers[device.ID]
	if found {
		return fmt.Errorf("workers already exists for device by ID: %d", device.ID)
	}

	return s.create(ctx, device)
}

func (s *WorkerStore) Update(ctx context.Context, device models.Conn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	worker, found := s.workers[device.ID]
	if !found {
		return fmt.Errorf("workers not found for device by ID: %d", device.ID)
	}
	if worker.conn.EQ(device) {
		return nil
	}

	for _, st := range worker.tokens {
		s.super.Remove(st)
	}

	return s.create(ctx, device)
}

func (s *WorkerStore) Delete(id int64) error {
	s.workersMu.Lock()
	worker, found := s.workers[id]
	if !found {
		s.workersMu.Unlock()
		return fmt.Errorf("workers not found for device by ID: %d", id)
	}

	for _, token := range worker.tokens {
		s.super.Remove(token)
	}
	delete(s.workers, id)
	s.workersMu.Unlock()
	return nil
}

func (s *WorkerStore) Register(bus *event.Bus) *WorkerStore {
	bus.OnDahuaDeviceCreated(func(ctx context.Context, evt event.DahuaDeviceCreated) error {
		return s.Create(ctx, evt.Conn)
	})
	bus.OnDahuaDeviceUpdated(func(ctx context.Context, evt event.DahuaDeviceUpdated) error {
		return s.Update(ctx, evt.Conn)
	})
	bus.OnDahuaDeviceDeleted(func(ctx context.Context, evt event.DahuaDeviceDeleted) error {
		return s.Delete(evt.DeviceID)
	})
	return s
}

func (s *WorkerStore) Bootstrap(ctx context.Context, db repo.DB, store *Store) error {
	ids, err := db.DahuaListDeviceIDs(ctx)
	if err != nil {
		return err
	}
	clients, err := store.ListClient(ctx, ids)
	if err != nil {
		return err
	}

	for _, conn := range clients {
		if err := s.Create(ctx, conn.Conn); err != nil {
			return err
		}
	}

	return err
}
