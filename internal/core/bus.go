package core

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func busLogErr(err error) {
	if err != nil {
		log.Err(err).Str("package", "dahuacore").Msg("Failed to handle event")
	}
}

func NewBus() *Bus {
	return &Bus{
		ServiceContext: sutureext.NewServiceContext("dahuacore.Bus"),
	}
}

type Bus struct {
	sutureext.ServiceContext
	onDahuaCameraCreated         []func(ctx context.Context, evt models.EventDahuaCameraCreated) error
	onDahuaCameraUpdated         []func(ctx context.Context, evt models.EventDahuaCameraUpdated) error
	onDahuaCameraDeleted         []func(ctx context.Context, evt models.EventDahuaCameraDeleted) error
	onDahuaCameraEvent           []func(ctx context.Context, evt models.EventDahuaCameraEvent) error
	onDahuaEventWorkerConnecting []func(ctx context.Context, evt models.EventDahuaEventWorkerConnecting) error
	onDahuaEventWorkerConnect    []func(ctx context.Context, evt models.EventDahuaEventWorkerConnect) error
	onDahuaEventWorkerDisconnect []func(ctx context.Context, evt models.EventDahuaEventWorkerDisconnect) error
}

func (dahuaBus *Bus) Register(pub pubsub.Pub) {
	dahuaBus.OnDahuaCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		return pub.Publish(ctx, evt)
	})
}

func (b *Bus) OnDahuaCameraCreated(h func(ctx context.Context, evt models.EventDahuaCameraCreated) error) {
	b.onDahuaCameraCreated = append(b.onDahuaCameraCreated, h)
}

func (b *Bus) OnDahuaCameraUpdated(h func(ctx context.Context, evt models.EventDahuaCameraUpdated) error) {
	b.onDahuaCameraUpdated = append(b.onDahuaCameraUpdated, h)
}

func (b *Bus) OnDahuaCameraDeleted(h func(ctx context.Context, evt models.EventDahuaCameraDeleted) error) {
	b.onDahuaCameraDeleted = append(b.onDahuaCameraDeleted, h)
}

func (b *Bus) OnDahuaCameraEvent(h func(ctx context.Context, evt models.EventDahuaCameraEvent) error) {
	b.onDahuaCameraEvent = append(b.onDahuaCameraEvent, h)
}

func (b *Bus) OnDahuaEventWorkerConnecting(h func(ctx context.Context, evt models.EventDahuaEventWorkerConnecting) error) {
	b.onDahuaEventWorkerConnecting = append(b.onDahuaEventWorkerConnecting, h)
}

func (b *Bus) OnDahuaEventWorkerConnect(h func(ctx context.Context, evt models.EventDahuaEventWorkerConnect) error) {
	b.onDahuaEventWorkerConnect = append(b.onDahuaEventWorkerConnect, h)
}

func (b *Bus) OnDahuaEventWorkerDisconnect(h func(ctx context.Context, evt models.EventDahuaEventWorkerDisconnect) error) {
	b.onDahuaEventWorkerDisconnect = append(b.onDahuaEventWorkerDisconnect, h)
}

func (b *Bus) DahuaCameraCreated(camera models.DahuaConn) {
	for _, v := range b.onDahuaCameraCreated {
		busLogErr(v(b.Context(), models.EventDahuaCameraCreated{
			Camera: camera,
		}))
	}
}

func (b *Bus) DahuaCameraUpdated(camera models.DahuaConn) {
	for _, v := range b.onDahuaCameraUpdated {
		busLogErr(v(b.Context(), models.EventDahuaCameraUpdated{
			Camera: camera,
		}))
	}
}

func (b *Bus) DahuaCameraDeleted(id int64) {
	for _, v := range b.onDahuaCameraDeleted {
		busLogErr(v(b.Context(), models.EventDahuaCameraDeleted{
			CameraID: id,
		}))
	}
}

func (b *Bus) DahuaCameraEvent(ctx context.Context, event models.DahuaEvent, eventRule models.DahuaEventRule) {
	for _, v := range b.onDahuaCameraEvent {
		busLogErr(v(ctx, models.EventDahuaCameraEvent{
			Event:     event,
			EventRule: eventRule,
		}))
	}
}

func (b *Bus) DahuaEventWorkerConnecting(cameraID int64) {
	for _, v := range b.onDahuaEventWorkerConnecting {
		busLogErr(v(b.Context(), models.EventDahuaEventWorkerConnecting{
			CameraID: cameraID,
		}))
	}
}

func (b *Bus) DahuaEventWorkerConnect(cameraID int64) {
	for _, v := range b.onDahuaEventWorkerConnect {
		busLogErr(v(b.Context(), models.EventDahuaEventWorkerConnect{
			CameraID: cameraID,
		}))
	}
}

func (b *Bus) DahuaEventWorkerDisconnect(cameraID int64, err error) {
	for _, v := range b.onDahuaEventWorkerDisconnect {
		busLogErr(v(b.Context(), models.EventDahuaEventWorkerDisconnect{
			CameraID: cameraID,
			Error:    err,
		}))
	}
}
