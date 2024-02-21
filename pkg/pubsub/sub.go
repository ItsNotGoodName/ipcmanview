package pubsub

import (
	"context"
)

type MiddlewareFunc func(next HandleFunc) HandleFunc

type Sub struct {
	pub   *Pub
	doneC <-chan struct{}
	errC  chan error
	id    int
}

// Wait blocks until the subscription is closed.
func (s Sub) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.doneC:
		err := <-s.errC
		s.errC <- err
		return err
	}
}

// Error should only be called after subscription is closed.
func (s Sub) Error() error {
	select {
	case <-s.doneC:
		err := <-s.errC
		s.errC <- err
		return err
	default:
		return nil
	}
}

func (s Sub) Close() error {
	s.pub.unsubscribe(s.id)
	return nil
}

type SubscribeBuilder struct {
	pub        *Pub
	topics     []string
	middleware []MiddlewareFunc
}

func (p *Pub) Subscribe(events ...Event) SubscribeBuilder {
	var topics []string
	for _, e := range events {
		topics = append(topics, e.EventTopic())
	}
	return SubscribeBuilder{
		pub:    p,
		topics: topics,
	}
}

// Middleware adds middleare between the handler.
func (b SubscribeBuilder) Middleware(fn MiddlewareFunc) SubscribeBuilder {
	b.middleware = append(b.middleware, fn)
	return b
}

// Function creates a subscription with a function.
func (b SubscribeBuilder) Function(ctx context.Context, fn HandleFunc) (Sub, error) {
	handle := fn
	for _, mw := range b.middleware {
		handle = mw(handle)
	}

	resC := make(chan int, 1)
	doneC := make(chan struct{})
	errC := make(chan error, 1)
	id := b.pub.subscribe(subscribe{
		topics: b.topics,
		handle: handle,
		resC:   resC,
		doneC:  doneC,
		errC:   errC,
	})

	return Sub{
		pub:   b.pub,
		doneC: doneC,
		errC:  errC,
		id:    id,
	}, nil
}

// Channel creates a subscription with a channel.
func (b SubscribeBuilder) Channel(ctx context.Context, size int) (Sub, <-chan Event, error) {
	eventsC := make(chan Event, size)

	sub, err := b.Function(ctx, func(ctx context.Context, event Event) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case eventsC <- event:
			return nil
		}
	})
	if err != nil {
		return Sub{}, nil, err
	}

	go func() {
		select {
		case <-ctx.Done():
			sub.Close()
			<-sub.doneC
		case <-sub.doneC:
		}
		// This assumes the publisher will not call the handle function after the subscription is fully closed
		close(eventsC)
	}()

	return sub, eventsC, nil
}
