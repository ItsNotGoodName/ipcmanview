package pubsub

import "context"

type Sub struct {
	doneC <-chan struct{}
	errC  chan error
	pub   Pub
	id    int
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
		}, nil
	}
}

// Wait blocks until the subscription is closed.
func (s *Sub) Wait(ctx context.Context) error {
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
func (s *Sub) Error() error {
	select {
	case <-s.pub.doneC:
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
	case <-ctx.Done():
		return ctx.Err()
	case <-s.pub.doneC:
		return ErrPubSubClosed
	case s.pub.commandC <- unsubscribe{id: s.id}:
		return nil
	}
}
