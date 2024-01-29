package pubsub

import (
	"context"
)

type Sub struct {
	doneC <-chan struct{}
	errC  chan error
	pub   Pub
	id    int
}

// Wait blocks until the subscription is closed.
func (s Sub) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.pub.doneC:
		return ErrPubSubClosed
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
	return s.close(context.Background())
}

func (s Sub) close(ctx context.Context) error {
	select {
	case <-s.doneC:
		return nil
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.pub.doneC:
		return ErrPubSubClosed
	case <-s.doneC:
		return nil
	case s.pub.unsubscribeC <- unsubscribe{id: s.id}:
		return nil
	}
}

type SubscriberBuilder struct {
	pub    Pub
	topics []string
}

func (p Pub) Subscribe(events ...Event) SubscriberBuilder {
	var topics []string
	for _, e := range events {
		topics = append(topics, e.EventTopic())
	}
	return SubscriberBuilder{
		pub:    p,
		topics: topics,
	}
}

func (b SubscriberBuilder) Function(ctx context.Context, handle HandleFunc) (Sub, error) {
	resC := make(chan int, 1)
	doneC := make(chan struct{})
	errC := make(chan error, 1)
	select {
	case <-ctx.Done():
		return Sub{}, ctx.Err()
	case <-b.pub.doneC:
		return Sub{}, ErrPubSubClosed
	case b.pub.subscribeC <- subscribe{
		topics: b.topics,
		handle: handle,
		resC:   resC,
		doneC:  doneC,
		errC:   errC,
	}:
	}

	select {
	case <-ctx.Done():
		return Sub{}, ctx.Err()
	case <-b.pub.doneC:
		return Sub{}, ErrPubSubClosed
	case id := <-resC:
		return Sub{
			doneC: doneC,
			pub:   b.pub,
			id:    id,
			errC:  errC,
		}, nil
	}
}

// Channel creates a subscription with a channel.
// The subscription is closed when the context is closed.
// The channel is closed when the subscription is closed.
func (b SubscriberBuilder) Channel(ctx context.Context, size int) (Sub, <-chan Event, error) {
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
		case <-b.pub.doneC:
		}
		// This assumes the publisher will not call the handle function after the subscription is fully closed
		close(eventsC)
	}()

	return sub, eventsC, nil
}
