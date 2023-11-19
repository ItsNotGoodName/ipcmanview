package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type EventHooks interface {
	CameraEvent(ctx context.Context, camera models.DahuaCamera, event models.DahuaEvent)
}

func newEventWorker(camera models.DahuaCamera, hooks EventHooks) eventWorker {
	return eventWorker{
		Camera: camera,
		hooks:  hooks,
	}
}

type eventWorker struct {
	Camera models.DahuaCamera
	hooks  EventHooks
}

func (w eventWorker) String() string {
	return fmt.Sprintf("dahua.eventWorker(id=%s)", w.Camera.ID)
}

func (w eventWorker) Serve(ctx context.Context) error {
	c := dahuacgi.NewConn(http.Client{}, NewAddress(w.Camera.Address), w.Camera.Username, w.Camera.Password)

	manager, err := dahuacgi.EventManagerGet(ctx, c, 0)
	if err != nil {
		var httpErr dahuacgi.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnauthorized {
			log.Err(err).Str("service", w.String()).Msg("Failed to get EventManager")
			return errors.Join(suture.ErrDoNotRestart, err)
		}

		return err
	}
	defer manager.Close()

	reader := manager.Reader()

	for {
		if err := reader.Poll(); err != nil {
			return err
		}

		rawEvent, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		event := NewDahuaEvent(rawEvent, time.Now())

		w.hooks.CameraEvent(ctx, w.Camera, event)
	}
}

type EventWorkerStore struct {
	super *suture.Supervisor
	hooks EventHooks

	workersMu sync.Mutex
	workers   map[string]suture.ServiceToken
}

func NewEventWorkerStore(super *suture.Supervisor, hooks EventHooks) *EventWorkerStore {
	return &EventWorkerStore{
		super:     super,
		hooks:     hooks,
		workersMu: sync.Mutex{},
		workers:   make(map[string]suture.ServiceToken),
	}
}

func (s *EventWorkerStore) Create(camera models.DahuaCamera) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	token, found := s.workers[camera.ID]
	if found {
		return fmt.Errorf("eventWorker already exists: %s", camera.ID)
	}

	log.Info().Str("id", camera.ID).Msg("Creating eventWorker")
	worker := newEventWorker(camera, s.hooks)
	token = s.super.Add(worker)
	s.workers[camera.ID] = token

	return nil
}

func (s *EventWorkerStore) Update(camera models.DahuaCamera) error {
	s.workersMu.Lock()
	defer s.workersMu.Unlock()

	token, found := s.workers[camera.ID]
	if !found {
		return fmt.Errorf("eventWorker not found by ID: %s", camera.ID)
	}

	log.Info().Str("id", camera.ID).Msg("Updating eventWorker")
	s.super.Remove(token)
	worker := newEventWorker(camera, s.hooks)
	token = s.super.Add(worker)
	s.workers[camera.ID] = token

	return nil
}

// func (s *EventSupervisor) Delete(id string) {
// 	s.workersMu.Lock()
// 	defer s.workersMu.Unlock()
//
// 	token, found := s.workers[id]
// 	if found {
// 		s.super.Remove(token)
// 	}
// 	delete(s.workers, id)
// }

type EventBus interface {
	OnCameraCreated(h func(evt models.EventDahuaCameraCreated) error)
	OnCameraUpdated(h func(evt models.EventDahuaCameraUpdated) error)
}

func RegisterEventBus(e *EventWorkerStore, bus EventBus) {
	bus.OnCameraCreated(func(evt models.EventDahuaCameraCreated) error {
		return e.Create(evt.Camera)
	})
	bus.OnCameraUpdated(func(evt models.EventDahuaCameraUpdated) error {
		return e.Update(evt.Camera)
	})
}
