package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type EventHooks interface {
	Connecting(ctx context.Context, deviceID int64)
	Connect(ctx context.Context, deviceID int64)
	Disconnect(deviceID int64, err error)
	Event(ctx context.Context, event models.DahuaEvent)
}

func NewEventWorker(device models.DahuaConn, hooks EventHooks) EventWorker {
	return EventWorker{
		device: device,
		hooks:  hooks,
	}
}

// EventWorker subscribes to events.
type EventWorker struct {
	device models.DahuaConn
	hooks  EventHooks
}

func (w EventWorker) String() string {
	return fmt.Sprintf("dahua.EventWorker(id=%d)", w.device.ID)
}

func (w EventWorker) Serve(ctx context.Context) error {
	w.hooks.Connecting(ctx, w.device.ID)
	err := w.serve(ctx)
	w.hooks.Disconnect(w.device.ID, err)
	return sutureext.SanitizeError(ctx, err)
}

func (w EventWorker) serve(ctx context.Context) error {
	c := dahuacgi.NewClient(http.Client{}, NewHTTPAddress(w.device.Address), w.device.Username, w.device.Password)

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

	w.hooks.Connect(ctx, w.device.ID)

	for reader := manager.Reader(); ; {
		if err := reader.Poll(); err != nil {
			return err
		}

		rawEvent, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		event := NewDahuaEvent(w.device.ID, rawEvent)

		w.hooks.Event(ctx, event)
	}
}
