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
	onCameraCreated []func(evt models.EventDahuaCameraCreated) error
	onCameraUpdated []func(evt models.EventDahuaCameraUpdated) error
	onCameraEvent   []func(ctx context.Context, evt models.EventDahuaCameraEvent) error
}

func (b *Bus) OnCameraCreated(h func(evt models.EventDahuaCameraCreated) error) {
	b.mu.Lock()
	b.onCameraCreated = append(b.onCameraCreated, h)
	b.mu.Unlock()
}

func (b *Bus) OnCameraUpdated(h func(evt models.EventDahuaCameraUpdated) error) {
	b.mu.Lock()
	b.onCameraUpdated = append(b.onCameraUpdated, h)
	b.mu.Unlock()
}

func (b *Bus) OnCameraEvent(h func(ctx context.Context, evt models.EventDahuaCameraEvent) error) {
	b.mu.Lock()
	b.onCameraEvent = append(b.onCameraEvent, h)
	b.mu.Unlock()
}

func (b *Bus) ConnCreated(conn Conn) {
	b.mu.Lock()
	for _, v := range b.onCameraCreated {
		busLogErr(v(models.EventDahuaCameraCreated{
			Camera: conn.Camera,
		}))
	}
	b.mu.Unlock()
}

func (b *Bus) ConnUpdated(conn Conn) {
	b.mu.Lock()
	for _, v := range b.onCameraUpdated {
		busLogErr(v(models.EventDahuaCameraUpdated{
			Camera: conn.Camera,
		}))
	}
	b.mu.Unlock()
}

func (b *Bus) CameraEvent(ctx context.Context, camera models.DahuaCamera, event models.DahuaEvent) {
	b.mu.Lock()
	for _, v := range b.onCameraEvent {
		busLogErr(v(ctx, models.EventDahuaCameraEvent{
			ID:    camera.ID,
			Event: event,
		}))
	}
	b.mu.Unlock()
}
