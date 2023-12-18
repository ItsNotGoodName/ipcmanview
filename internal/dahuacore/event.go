package dahuacore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

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

func NewEventWorker(camera models.DahuaConn, hooks EventHooks) EventWorker {
	return EventWorker{
		Camera: camera,
		hooks:  hooks,
	}
}

type EventWorker struct {
	Camera models.DahuaConn
	hooks  EventHooks
}

func (w EventWorker) String() string {
	return fmt.Sprintf("dahuacore.EventWorker(id=%d)", w.Camera.ID)
}

func (w EventWorker) Serve(ctx context.Context) error {
	w.hooks.Connecting(ctx, w.Camera.ID)
	err := w.serve(ctx)
	w.hooks.Disconnect(w.Camera.ID, err)
	return err
}

func (w EventWorker) serve(ctx context.Context) error {
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
