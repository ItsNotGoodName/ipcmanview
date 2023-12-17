package dahuacore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type EventHooks interface {
	Connecting(ctx context.Context, cameraID int64)
	Connect(ctx context.Context, cameraID int64)
	Disconnect(cameraID int64, err error)
	Event(ctx context.Context, event models.DahuaEvent)
}

func newEventWorker(camera models.DahuaConn, hooks EventHooks) eventWorker {
	return eventWorker{
		Camera: camera,
		hooks:  hooks,
	}
}

type eventWorker struct {
	Camera models.DahuaConn
	hooks  EventHooks
}

func (w eventWorker) String() string {
	return fmt.Sprintf("dahuacore.eventWorker(id=%d)", w.Camera.ID)
}

func (w eventWorker) Serve(ctx context.Context) error {
	w.hooks.Connecting(ctx, w.Camera.ID)
	err := w.serve(ctx)
	w.hooks.Disconnect(w.Camera.ID, err)
	return err
}

func (w eventWorker) serve(ctx context.Context) error {
	c := dahuacgi.NewClient(http.Client{}, NewHTTPAddress(w.Camera.Address), w.Camera.Username, w.Camera.Password)

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

	w.hooks.Connect(ctx, w.Camera.ID)

	for reader := manager.Reader(); ; {
		if err := reader.Poll(); err != nil {
			return err
		}

		rawEvent, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		event := NewDahuaEvent(w.Camera.ID, rawEvent)

		w.hooks.Event(ctx, event)
	}
}

type EventWorkerStore struct {
	super *suture.Supervisor
	hooks EventHooks

	workersMu sync.Mutex
	workers   map[int64]suture.ServiceToken
}

func NewEventWorkerStore(super *suture.Supervisor, hooks EventHooks) *EventWorkerStore {
	return &EventWorkerStore{
		super:     super,
		hooks:     hooks,
		workersMu: sync.Mutex{},
		workers:   make(map[int64]suture.ServiceToken),
	}
}

func (s *EventWorkerStore) Create(camera models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	token, found := s.workers[camera.ID]
	if found {
		return fmt.Errorf("eventWorker already exists: %d", camera.ID)
	}

	log.Info().Int64("id", camera.ID).Msg("Creating dahua.eventWorker")
	worker := newEventWorker(camera, s.hooks)
	token = s.super.Add(worker)
	s.workers[camera.ID] = token

	return nil
}

func (s *EventWorkerStore) Update(camera models.DahuaConn) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	token, found := s.workers[camera.ID]
	if !found {
		return fmt.Errorf("eventWorker not found by ID: %d", camera.ID)
	}

	log.Info().Int64("id", camera.ID).Msg("Updating eventWorker")
	s.super.Remove(token)
	worker := newEventWorker(camera, s.hooks)
	token = s.super.Add(worker)
	s.workers[camera.ID] = token

	return nil
}

func (s *EventWorkerStore) Delete(id int64) {
	s.workersMu.Lock()
	token, found := s.workers[id]
	if found {
		s.super.Remove(token)
	}
	delete(s.workers, id)
	s.workersMu.Unlock()
}

func (e *EventWorkerStore) Register(bus *core.Bus) {
	bus.OnEventDahuaCameraCreated(func(ctx context.Context, evt models.EventDahuaCameraCreated) error {
		return e.Create(evt.Camera)
	})
	bus.OnEventDahuaCameraUpdated(func(ctx context.Context, evt models.EventDahuaCameraUpdated) error {
		return e.Update(evt.Camera)
	})
	bus.OnEventDahuaCameraDeleted(func(ctx context.Context, evt models.EventDahuaCameraDeleted) error {
		e.Delete(evt.CameraID)
		return nil
	})
}
