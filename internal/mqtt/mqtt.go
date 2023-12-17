package mqtt

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

const DahuaEventType = "dahua_event"

type DahuaEvent struct {
	models.DahuaEvent
	// EventType is required for HomeAssistant.
	EventType string `json:"event_type"`
}

func Register(conn Conn, bus *core.Bus) error {
	bus.OnEventDahuaCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		conn.Ready()

		if evt.EventRule.IgnoreMQTT {
			return nil
		}

		b, err := json.Marshal(DahuaEvent{DahuaEvent: evt.Event, EventType: DahuaEventType})
		if err != nil {
			return err
		}

		return Wait(conn.Client.Publish(conn.Topic.Join("dahua", strconv.FormatInt(evt.Event.CameraID, 10), "event"), 0, false, b))
	})
	bus.OnEventDahuaEventWorkerConnect(func(ctx context.Context, evt models.EventDahuaEventWorkerConnect) error {
		conn.Ready()

		if err := publishDahuaEventError(ctx, conn, evt.CameraID, nil); err != nil {
			return err
		}

		return Wait(conn.Client.Publish(conn.Topic.Join("dahua", strconv.FormatInt(evt.CameraID, 10), "event", "state"), 0, true, "online"))
	})
	bus.OnEventDahuaEventWorkerDisconnect(func(ctx context.Context, evt models.EventDahuaEventWorkerDisconnect) error {
		conn.Ready()

		if err := publishDahuaEventError(ctx, conn, evt.CameraID, evt.Error); err != nil {
			return err
		}

		return Wait(conn.Client.Publish(conn.Topic.Join("dahua", strconv.FormatInt(evt.CameraID, 10), "event", "state"), 0, true, "offline"))
	})
	return nil
}

func publishDahuaEventError(ctx context.Context, conn Conn, cameraID int64, err error) error {
	var payload any
	if err != nil {
		payload = err.Error()
	} else {
		payload = []byte{}
	}
	return Wait(conn.Client.Publish(conn.Topic.Join("dahua", strconv.FormatInt(cameraID, 10), "event", "error"), 0, true, payload))
}
