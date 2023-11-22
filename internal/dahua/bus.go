package dahua

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/rs/zerolog/log"
)

func busLogErr(err error) {
	if err != nil {
		log.Err(err).Caller().Send()
	}
}

func NewBus() *Bus {
	return &Bus{}
}

type Bus struct {
	mu              sync.Mutex
	onCameraCreated []func(ctx context.Context, evt models.EventDahuaCameraCreated) error
	onCameraUpdated []func(ctx context.Context, evt models.EventDahuaCameraUpdated) error
	onCameraDeleted []func(ctx context.Context, evt models.EventDahuaCameraDeleted) error
	onCameraEvent   []func(ctx context.Context, evt models.EventDahuaCameraEvent) error
}

func (b *Bus) OnCameraCreated(h func(ctx context.Context, evt models.EventDahuaCameraCreated) error) {
	b.mu.Lock()
	b.onCameraCreated = append(b.onCameraCreated, h)
	b.mu.Unlock()
}

func (b *Bus) OnCameraUpdated(h func(ctx context.Context, evt models.EventDahuaCameraUpdated) error) {
	b.mu.Lock()
	b.onCameraUpdated = append(b.onCameraUpdated, h)
	b.mu.Unlock()
}

func (b *Bus) OnCameraDeleted(h func(ctx context.Context, evt models.EventDahuaCameraDeleted) error) {
	b.mu.Lock()
	b.onCameraDeleted = append(b.onCameraDeleted, h)
	b.mu.Unlock()
}

func (b *Bus) OnCameraEvent(h func(ctx context.Context, evt models.EventDahuaCameraEvent) error) {
	b.mu.Lock()
	b.onCameraEvent = append(b.onCameraEvent, h)
	b.mu.Unlock()
}

func (b *Bus) CameraCreated(camera models.DahuaCamera) {
	b.mu.Lock()
	for _, v := range b.onCameraCreated {
		busLogErr(v(context.TODO(), models.EventDahuaCameraCreated{
			Camera: camera,
		}))
	}
	b.mu.Unlock()
}

func (b *Bus) CameraUpdated(camera models.DahuaCamera) {
	b.mu.Lock()
	for _, v := range b.onCameraUpdated {
		busLogErr(v(context.TODO(), models.EventDahuaCameraUpdated{
			Camera: camera,
		}))
	}
	b.mu.Unlock()
}

func (b *Bus) CameraDeleted(id int64) {
	b.mu.Lock()
	for _, v := range b.onCameraDeleted {
		busLogErr(v(context.TODO(), models.EventDahuaCameraDeleted{
			CameraID: id,
		}))
	}
	b.mu.Unlock()
}

func (b *Bus) CameraEvent(event models.DahuaEvent) {
	b.mu.Lock()
	for _, v := range b.onCameraEvent {
		busLogErr(v(context.TODO(), models.EventDahuaCameraEvent{
			Event: event,
		}))
	}
	b.mu.Unlock()
}
