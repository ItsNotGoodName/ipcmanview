package dahua

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/thejerf/suture/v4"
)

func DefaultWorkerBuilder(hooks DefaultEventHooks, bus *core.Bus, store *Store, conn ConnRepo) WorkerBuilder {
	return func(ctx context.Context, super *suture.Supervisor, device models.DahuaConn) ([]suture.ServiceToken, error) {
		var tokens []suture.ServiceToken

		{
			worker := NewEventWorker(device, hooks)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		{
			worker := NewCoaxialWorker(bus, device.ID, store, conn)
			token := super.Add(worker)
			tokens = append(tokens, token)
		}

		return tokens, nil
	}
}

type WorkerBuilder = func(ctx context.Context, super *suture.Supervisor, device models.DahuaConn) ([]suture.ServiceToken, error)

type WorkerStore struct {
	super *suture.Supervisor
	build WorkerBuilder

	workersMu sync.Mutex
	workers   map[int64]workerData
}

type workerData struct {
	device models.DahuaConn
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

func (s *WorkerStore) create(ctx context.Context, device models.DahuaConn) error {
	tokens, err := s.build(ctx, s.super, device)
	if err != nil {
		return err
	}

	s.workers[device.ID] = workerData{
		device: device,
		tokens: tokens,
	}

	return nil
}

func (s *WorkerStore) Create(ctx context.Context, device models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	_, found := s.workers[device.ID]
	if found {
		return fmt.Errorf("workers already exists for device by ID: %d", device.ID)
	}

	return s.create(ctx, device)
}

func (s *WorkerStore) Update(ctx context.Context, device models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	worker, found := s.workers[device.ID]
	if !found {
		return fmt.Errorf("workers not found for device by ID: %d", device.ID)
	}
	if ConnEqual(worker.device, device) {
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

func (e *WorkerStore) Register(bus *core.Bus) {
	bus.OnEventDahuaDeviceCreated(func(ctx context.Context, evt models.EventDahuaDeviceCreated) error {
		return e.Create(ctx, evt.Device.DahuaConn)
	})
	bus.OnEventDahuaDeviceUpdated(func(ctx context.Context, evt models.EventDahuaDeviceUpdated) error {
		return e.Update(ctx, evt.Device.DahuaConn)
	})
	bus.OnEventDahuaDeviceDeleted(func(ctx context.Context, evt models.EventDahuaDeviceDeleted) error {
		return e.Delete(evt.DeviceID)
	})
}

type DeviceStore interface {
	ListConn(ctx context.Context) ([]models.DahuaConn, error)
}

func (w *WorkerStore) Bootstrap(ctx context.Context, deviceStore DeviceStore, store *Store) error {
	devices, err := deviceStore.ListConn(ctx)
	if err != nil {
		return err
	}
	conns := store.ConnList(ctx, devices)
	for _, conn := range conns {
		if err := w.Create(ctx, conn.Conn); err != nil {
			return err
		}
	}
	return err
}
