package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuaevents"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type WorkerFactory = func(ctx context.Context, super *suture.Supervisor, device Conn) []suture.ServiceToken

// WorkerManager manages devices workers.
type WorkerManager struct {
	super   *suture.Supervisor
	factory WorkerFactory

	workersMu sync.Mutex
	workers   map[int64]workerData
}

func (*WorkerManager) String() string {
	return "dahua.WorkerManager"
}

type workerData struct {
	conn   Conn
	tokens []suture.ServiceToken
}

func NewWorkerManager(super *suture.Supervisor, factory WorkerFactory) *WorkerManager {
	return &WorkerManager{
		super:     super,
		factory:   factory,
		workersMu: sync.Mutex{},
		workers:   make(map[int64]workerData),
	}
}

func (m *WorkerManager) Upsert(ctx context.Context, conn Conn) error {
	m.workersMu.Lock()
	defer m.workersMu.Unlock()

	worker, found := m.workers[conn.ID]
	if found {
		if worker.conn.EQ(conn) {
			return nil
		}

		for _, st := range worker.tokens {
			m.super.Remove(st)
		}
	}

	m.workers[conn.ID] = workerData{
		conn:   conn,
		tokens: m.factory(ctx, m.super, conn),
	}

	return nil
}

func (m *WorkerManager) Delete(id int64) error {
	m.workersMu.Lock()
	worker, found := m.workers[id]
	if !found {
		m.workersMu.Unlock()
		return nil
	}

	for _, token := range worker.tokens {
		m.super.Remove(token)
	}
	delete(m.workers, id)
	m.workersMu.Unlock()
	return nil
}

func (m *WorkerManager) Register(bus *event.Bus, db sqlite.DB) *WorkerManager {
	upsert := func(ctx context.Context, deviceID int64) error {
		conn, err := GetConn(ctx, db, deviceID)
		if err != nil {
			if core.IsNotFound(err) {
				return m.Delete(deviceID)
			}
			return err
		}

		return m.Upsert(ctx, conn)
	}

	bus.OnDahuaDeviceCreated(m.String(), func(ctx context.Context, evt event.DahuaDeviceCreated) error {
		return upsert(ctx, evt.DeviceID)
	})
	bus.OnDahuaDeviceUpdated(m.String(), func(ctx context.Context, evt event.DahuaDeviceUpdated) error {
		return upsert(ctx, evt.DeviceID)
	})
	bus.OnDahuaDeviceDeleted(m.String(), func(ctx context.Context, evt event.DahuaDeviceDeleted) error {
		return m.Delete(evt.DeviceID)
	})

	return m
}

func (m *WorkerManager) Bootstrap(ctx context.Context, db sqlite.DB, store *Store) error {
	conns, err := ListConn(ctx, db)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := m.Upsert(ctx, conn); err != nil {
			return err
		}
	}

	return err
}

type Worker struct {
	DeviceID int64
	Type     models.DahuaWorkerType
}

type WorkerHooks interface {
	Serve(ctx context.Context, w Worker, connected bool, fn func(ctx context.Context) error) error
	Connected(ctx context.Context, w Worker)
}

func NewQuickScanWorker(hooks WorkerHooks, pub *pubsub.Pub, bus *event.Bus, db sqlite.DB, store *Store, scanLockStore ScanLockStore, deviceID int64) QuickScanWorker {
	return QuickScanWorker{
		hooks:         hooks,
		worker:        Worker{DeviceID: deviceID, Type: models.DahuaWorkerTypeQuickScan},
		pub:           pub,
		bus:           bus,
		db:            db,
		store:         store,
		scanLockStore: scanLockStore,
		deviceID:      deviceID,
	}
}

// QuickScanWorker scans devices for files.
type QuickScanWorker struct {
	hooks         WorkerHooks
	worker        Worker
	pub           *pubsub.Pub
	bus           *event.Bus
	db            sqlite.DB
	store         *Store
	scanLockStore ScanLockStore
	deviceID      int64
}

func (w QuickScanWorker) String() string {
	return fmt.Sprintf("dahua.QuickScanWorker(id=%d)", w.deviceID)
}

func (w QuickScanWorker) Serve(ctx context.Context) error {
	err := w.hooks.Serve(ctx, w.worker, true, w.serve)
	return sutureext.SanitizeError(ctx, err)
}

func (w QuickScanWorker) serve(ctx context.Context) error {
	// New file was created
	newFileC := make(chan struct{}, 1)

	// Subscribe
	sub, err := w.pub.
		Subscribe().
		Function(func(ctx context.Context, evt pubsub.Event) error {
			switch e := evt.(type) {
			case event.DahuaEvent:
				if e.Event.DeviceID != w.deviceID {
					return nil
				}

				switch e.Event.Code {
				case dahuaevents.CodeNewFile:
					core.FlagChannel(newFileC)
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	defer sub.Close()

	// Scan after a duration
	timer := time.NewTimer(0)
	const timerDuration = 10 * time.Second
	var timerEnd time.Time
	timerActive := false

	// Scan every 30 minutes
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case fired := <-timer.C:
			if err := w.scan(ctx); err != nil {
				return err
			}

			if fired.Before(timerEnd) {
				// Start timer
				timer.Reset(timerDuration)
				timerActive = true
			} else {
				timerActive = false
			}
		case <-ticker.C:
			// Start timer
			timer.Reset(0)
			timerActive = true
		case <-newFileC:
			timerEnd = time.Now().Add(timerDuration / 2)
			if timerActive {
				continue
			}

			// Start timer
			timer.Reset(timerDuration)
			timerActive = true
		}
	}
}

func (w QuickScanWorker) scan(ctx context.Context) error {
	unlock, err := w.scanLockStore.Lock(ctx, w.deviceID)
	if err != nil {
		return err
	}
	defer unlock()

	client, err := w.store.GetClient(ctx, w.deviceID)
	if err != nil {
		return err
	}

	return scan(ctx, w.db, w.bus, client.RPC, client.Conn, models.DahuaScanTypeQuick)
}

func NewEventWorker(hooks WorkerHooks, db sqlite.DB, bus *event.Bus, conn Conn) EventWorker {
	return EventWorker{
		hooks: hooks,
		worker: Worker{
			DeviceID: conn.ID,
			Type:     models.DahuaWorkerTypeEvent,
		},
		db:     db,
		bus:    bus,
		device: conn,
	}
}

// EventWorker publishes events to bus.
type EventWorker struct {
	hooks  WorkerHooks
	worker Worker
	db     sqlite.DB
	bus    *event.Bus
	device Conn
}

func (w EventWorker) String() string {
	return fmt.Sprintf("dahua.EventWorker(id=%d)", w.device.ID)
}

func (w EventWorker) Serve(ctx context.Context) error {
	err := w.hooks.Serve(ctx, w.worker, false, w.serve)
	return sutureext.SanitizeError(ctx, err)
}

func (w EventWorker) serve(ctx context.Context) error {
	c := dahuacgi.NewClient(http.Client{}, w.device.URL, w.device.Username, w.device.Password)

	manager, err := dahuacgi.EventManagerGet(ctx, c, 0)
	if err != nil {
		var httpErr dahuacgi.HTTPError
		if errors.As(err, &httpErr) && slices.Contains([]int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
		}, httpErr.StatusCode) {
			log.Err(err).Str("service", w.String()).Msg("Failed to get EventManager")
			return errors.Join(suture.ErrDoNotRestart, err)
		}

		return err
	}
	defer manager.Close()

	w.hooks.Connected(ctx, w.worker)

	for reader := manager.Reader(); ; {
		if err := reader.Poll(); err != nil {
			return err
		}

		evt, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		if err = publishEvent(ctx, w.db, w.bus, w.device.ID, evt); err != nil {
			return err
		}
	}
}

func NewCoaxialWorker(hooks WorkerHooks, db sqlite.DB, bus *event.Bus, store *Store, deviceID int64) CoaxialWorker {
	return CoaxialWorker{
		hooks: hooks,
		worker: Worker{
			DeviceID: deviceID,
			Type:     models.DahuaWorkerTypeCoaxial,
		},
		db:       db,
		bus:      bus,
		store:    store,
		deviceID: deviceID,
	}
}

// CoaxialWorker publishes coaxial status to the bus.
type CoaxialWorker struct {
	hooks    WorkerHooks
	worker   Worker
	db       sqlite.DB
	bus      *event.Bus
	store    *Store
	deviceID int64
}

func (w CoaxialWorker) String() string {
	return fmt.Sprintf("dahua.CoaxialWorker(id=%d)", w.deviceID)
}

func (w CoaxialWorker) Serve(ctx context.Context) error {
	err := w.hooks.Serve(ctx, w.worker, true, w.serve)
	return sutureext.SanitizeError(ctx, err)
}

func (w CoaxialWorker) serve(ctx context.Context) error {
	client, err := w.store.GetClient(ctx, w.deviceID)
	if err != nil {
		return err
	}

	channel := 1

	// Does this device support coaxial?
	caps, err := GetCoaxialCaps(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	if !(caps.SupportControlSpeaker || caps.SupportControlLight || caps.SupportControlFullcolorLight) {
		return suture.ErrDoNotRestart
	}

	// Get and publish initial coaxial status
	coaxialStatus, err := GetCoaxialStatus(ctx, client.RPC, channel)
	if err != nil {
		return err
	}
	w.bus.DahuaCoaxialStatus(event.DahuaCoaxialStatus{
		DeviceID:      w.deviceID,
		Channel:       channel,
		CoaxialStatus: coaxialStatus,
	})

	t := time.NewTicker(1 * time.Second)

	// Get and send coaxial status if it changes on an interval
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}

		s, err := GetCoaxialStatus(ctx, client.RPC, channel)
		if err != nil {
			return err
		}
		if coaxialStatus.Speaker == s.Speaker && coaxialStatus.WhiteLight == s.WhiteLight {
			continue
		}
		coaxialStatus = s

		w.bus.DahuaCoaxialStatus(event.DahuaCoaxialStatus{
			DeviceID:      w.deviceID,
			Channel:       channel,
			CoaxialStatus: coaxialStatus,
		})
	}
}
