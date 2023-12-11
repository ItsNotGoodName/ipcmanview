package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/rs/zerolog/log"
)

func busLogErr(err error) {
	if err != nil {
		log.Err(err).Msg("Failed to handle event")
	}
}

func NewBus() *Bus {
	return &Bus{}
}

type Bus struct {
	onCameraCreated []func(ctx context.Context, evt models.EventDahuaCameraCreated) error
	onCameraUpdated []func(ctx context.Context, evt models.EventDahuaCameraUpdated) error
	onCameraDeleted []func(ctx context.Context, evt models.EventDahuaCameraDeleted) error
	onCameraEvent   []func(ctx context.Context, evt models.EventDahuaCameraEvent) error
}

func (dahuaBus *Bus) Register(pub pubsub.Pub) {
	dahuaBus.OnCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		return pub.Publish(ctx, evt)
	})
}

func (b *Bus) OnCameraCreated(h func(ctx context.Context, evt models.EventDahuaCameraCreated) error) {
	b.onCameraCreated = append(b.onCameraCreated, h)
}

func (b *Bus) OnCameraUpdated(h func(ctx context.Context, evt models.EventDahuaCameraUpdated) error) {
	b.onCameraUpdated = append(b.onCameraUpdated, h)
}

func (b *Bus) OnCameraDeleted(h func(ctx context.Context, evt models.EventDahuaCameraDeleted) error) {
	b.onCameraDeleted = append(b.onCameraDeleted, h)
}

func (b *Bus) OnCameraEvent(h func(ctx context.Context, evt models.EventDahuaCameraEvent) error) {
	b.onCameraEvent = append(b.onCameraEvent, h)
}

func (b *Bus) CameraCreated(camera models.DahuaConn) {
	for _, v := range b.onCameraCreated {
		busLogErr(v(context.TODO(), models.EventDahuaCameraCreated{
			Camera: camera,
		}))
	}
}

func (b *Bus) CameraUpdated(camera models.DahuaConn) {
	for _, v := range b.onCameraUpdated {
		busLogErr(v(context.TODO(), models.EventDahuaCameraUpdated{
			Camera: camera,
		}))
	}
}

func (b *Bus) CameraDeleted(id int64) {
	for _, v := range b.onCameraDeleted {
		busLogErr(v(context.TODO(), models.EventDahuaCameraDeleted{
			CameraID: id,
		}))
	}
}

func (b *Bus) CameraEvent(ctx context.Context, event models.DahuaEvent) {
	for _, v := range b.onCameraEvent {
		busLogErr(v(context.TODO(), models.EventDahuaCameraEvent{
			Event: event,
		}))
	}
}
