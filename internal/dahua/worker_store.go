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
	db      repo.DB
	factory WorkerFactory

	workersMu sync.Mutex
	workers   map[int64]workerData
}

type workerData struct {
	conn   models.Conn
	tokens []suture.ServiceToken
}

func NewWorkerStore(super *suture.Supervisor, db repo.DB, factory WorkerFactory) *WorkerStore {
	return &WorkerStore{
		super:     super,
		db:        db,
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

func (s *WorkerStore) Create(ctx context.Context, deviceID int64) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	conn, err := s.db.DahuaGetConn(ctx, deviceID)
	if err != nil {
		return err
	}

	_, found := s.workers[deviceID]
	if found {
		return fmt.Errorf("workers already exists for device by ID: %d", deviceID)
	}

	return s.create(ctx, conn)
}

func (s *WorkerStore) Update(ctx context.Context, deviceID int64) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	conn, err := s.db.DahuaGetConn(ctx, deviceID)
	if err != nil {
		return err
	}

	worker, found := s.workers[deviceID]
	if !found {
		return fmt.Errorf("workers not found for device by ID: %d", deviceID)
	}
	if worker.conn.EQ(conn) {
		return nil
	}

	for _, st := range worker.tokens {
		s.super.Remove(st)
	}

	return s.create(ctx, conn)
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
		return s.Create(ctx, evt.DeviceID)
	})
	bus.OnDahuaDeviceUpdated(func(ctx context.Context, evt event.DahuaDeviceUpdated) error {
		return s.Update(ctx, evt.DeviceID)
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

	for _, id := range ids {
		if err := s.Create(ctx, id); err != nil {
			if repo.IsNotFound(err) {
				continue
			}
			return err
		}
	}

	return err
}
