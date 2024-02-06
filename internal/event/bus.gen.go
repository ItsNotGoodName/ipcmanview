// Code generated by generate-bus.go; DO NOT EDIT.
package event

import (
	"context"
	"errors"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func busLogError(err error) {
	if err != nil {
		log.Err(err).Str("package", "event").Send()
	}
}

func NewBus() *Bus {
	return &Bus{
		ServiceContext: sutureext.NewServiceContext("event.Bus"),
	}
}

type Bus struct {
	sutureext.ServiceContext
	onEventQueued []func(ctx context.Context, event EventQueued) error
	onEvent []func(ctx context.Context, event Event) error
	onDahuaEvent []func(ctx context.Context, event DahuaEvent) error
	onDahuaEventWorkerConnecting []func(ctx context.Context, event DahuaEventWorkerConnecting) error
	onDahuaEventWorkerConnect []func(ctx context.Context, event DahuaEventWorkerConnect) error
	onDahuaEventWorkerDisconnect []func(ctx context.Context, event DahuaEventWorkerDisconnect) error
	onDahuaCoaxialStatus []func(ctx context.Context, event DahuaCoaxialStatus) error
}

func (b *Bus) Register(pub pubsub.Pub) (*Bus) {
	b.OnEventQueued(func(ctx context.Context, evt EventQueued) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnEvent(func(ctx context.Context, evt Event) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEvent(func(ctx context.Context, evt DahuaEvent) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEventWorkerConnecting(func(ctx context.Context, evt DahuaEventWorkerConnecting) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEventWorkerConnect(func(ctx context.Context, evt DahuaEventWorkerConnect) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaEventWorkerDisconnect(func(ctx context.Context, evt DahuaEventWorkerDisconnect) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	b.OnDahuaCoaxialStatus(func(ctx context.Context, evt DahuaCoaxialStatus) error {
		err := pub.Publish(ctx, evt)
		if err == nil || errors.Is(err, pubsub.ErrPubSubClosed) {
			return nil
		}
		return err
	})
	return b
}


func (b *Bus) OnEventQueued(h func(ctx context.Context, evt EventQueued) error) {
	b.onEventQueued = append(b.onEventQueued, h)
}

func (b *Bus) OnEvent(h func(ctx context.Context, evt Event) error) {
	b.onEvent = append(b.onEvent, h)
}

func (b *Bus) OnDahuaEvent(h func(ctx context.Context, evt DahuaEvent) error) {
	b.onDahuaEvent = append(b.onDahuaEvent, h)
}

func (b *Bus) OnDahuaEventWorkerConnecting(h func(ctx context.Context, evt DahuaEventWorkerConnecting) error) {
	b.onDahuaEventWorkerConnecting = append(b.onDahuaEventWorkerConnecting, h)
}

func (b *Bus) OnDahuaEventWorkerConnect(h func(ctx context.Context, evt DahuaEventWorkerConnect) error) {
	b.onDahuaEventWorkerConnect = append(b.onDahuaEventWorkerConnect, h)
}

func (b *Bus) OnDahuaEventWorkerDisconnect(h func(ctx context.Context, evt DahuaEventWorkerDisconnect) error) {
	b.onDahuaEventWorkerDisconnect = append(b.onDahuaEventWorkerDisconnect, h)
}

func (b *Bus) OnDahuaCoaxialStatus(h func(ctx context.Context, evt DahuaCoaxialStatus) error) {
	b.onDahuaCoaxialStatus = append(b.onDahuaCoaxialStatus, h)
}



func (b *Bus) EventQueued(evt EventQueued) {
	for _, v := range b.onEventQueued {
		busLogError(v(b.Context(), evt))
	}
}

func (b *Bus) Event(evt Event) {
	for _, v := range b.onEvent {
		busLogError(v(b.Context(), evt))
	}
}

func (b *Bus) DahuaEvent(evt DahuaEvent) {
	for _, v := range b.onDahuaEvent {
		busLogError(v(b.Context(), evt))
	}
}

func (b *Bus) DahuaEventWorkerConnecting(evt DahuaEventWorkerConnecting) {
	for _, v := range b.onDahuaEventWorkerConnecting {
		busLogError(v(b.Context(), evt))
	}
}

func (b *Bus) DahuaEventWorkerConnect(evt DahuaEventWorkerConnect) {
	for _, v := range b.onDahuaEventWorkerConnect {
		busLogError(v(b.Context(), evt))
	}
}

func (b *Bus) DahuaEventWorkerDisconnect(evt DahuaEventWorkerDisconnect) {
	for _, v := range b.onDahuaEventWorkerDisconnect {
		busLogError(v(b.Context(), evt))
	}
}

func (b *Bus) DahuaCoaxialStatus(evt DahuaCoaxialStatus) {
	for _, v := range b.onDahuaCoaxialStatus {
		busLogError(v(b.Context(), evt))
	}
}

