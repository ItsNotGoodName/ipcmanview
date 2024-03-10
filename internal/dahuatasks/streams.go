package dahuatasks

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/squeuel"
)

type StreamPayload struct {
	DeviceID int64
}

func (p StreamPayload) TaskID() squeuel.Option {
	return squeuel.TaskID(fmt.Sprintf("%d", p.DeviceID))
}

var (
	SyncStreamTask = squeuel.NewTaskBuilder[StreamPayload]("dahua-stream:sync")
	PushStreamTask = squeuel.NewTaskBuilder[StreamPayload]("dahua-stream:push")
)

func RegisterStreams() {
	enqueue := func(ctx context.Context, deviceID int64) error {
		payload := StreamPayload{
			DeviceID: deviceID,
		}
		task, err := SyncStreamTask.New(payload, payload.TaskID())
		if err != nil {
			return err
		}

		_, err = squeuel.EnqueueTask(ctx, app.DB, app.Hub, task)
		return err
	}
	app.Hub.OnDahuaDeviceCreated("dahua.SyncStreams", func(ctx context.Context, event bus.DahuaDeviceCreated) error {
		return enqueue(ctx, event.DeviceID)
	})
	app.Hub.OnDahuaDeviceUpdated("dahua.SyncStreams", func(ctx context.Context, event bus.DahuaDeviceUpdated) error {
		return enqueue(ctx, event.DeviceID)
	})
}

func HandleSyncStreamTask(ctx context.Context, task *squeuel.Task) error {
	payload, err := SyncStreamTask.Payload(task)
	if err != nil {
		return err
	}

	conn, err := dahua.GetClient(ctx, payload.DeviceID)
	if err != nil {
		return err
	}

	if !dahua.SupportStream(conn.Conn.Feature) {
		return nil
	}

	if err := dahua.SyncStreams(ctx, payload.DeviceID, conn.RPC); err != nil {
		return err
	}

	task, err = PushStreamTask.New(payload, payload.TaskID())
	if err != nil {
		return err
	}

	_, err = squeuel.EnqueueTask(ctx, app.DB, app.Hub, task)
	return err
}

func HandlePushStreamTask(ctx context.Context, task *squeuel.Task) error {
	payload, err := PushStreamTask.Payload(task)
	if err != nil {
		return err
	}

	return dahua.PushStreams(ctx, payload.DeviceID)
}
