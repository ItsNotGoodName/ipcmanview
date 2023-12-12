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
	case s.pub.commandC <- unsubscribe{id: s.id}:
		return nil
	}
}

func (p Pub) Subscribe(ctx context.Context, handle HandleFunc, events ...Event) (Sub, error) {
	var topics []string
	for _, e := range events {
		topics = append(topics, e.EventName())
	}

	resC := make(chan int, 1)
	doneC := make(chan struct{})
	errC := make(chan error, 1)
	select {
	case <-ctx.Done():
		return Sub{}, ctx.Err()
	case <-p.doneC:
		return Sub{}, ErrPubSubClosed
	case p.commandC <- subscribe{
		topics: topics,
		handle: handle,
		resC:   resC,
		doneC:  doneC,
		errC:   errC,
	}:
	}

	select {
	case <-ctx.Done():
		return Sub{}, ctx.Err()
	case <-p.doneC:
		return Sub{}, ErrPubSubClosed
	case id := <-resC:
		return Sub{
			doneC: doneC,
			pub:   p,
			id:    id,
			errC:  errC,
		}, nil
	}
}

// SubscribeChan creates a subscription with a channel.
// The subscription is closed when the context is closed.
// The channel is closed when the subscription is closed.
func (p Pub) SubscribeChan(ctx context.Context, size int, events ...Event) (Sub, <-chan Event, error) {
	eventsC := make(chan Event, size)

	sub, err := p.Subscribe(ctx, func(ctx context.Context, event Event) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case eventsC <- event:
			return nil
		}
	}, events...)
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
