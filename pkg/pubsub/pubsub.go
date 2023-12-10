// pubsub is a simple in-memory event pub/sub.
package pubsub

import (
	"context"
	"errors"
	"slices"
	"time"
)

var ErrPubSubClosed = errors.New("pub sub closed")

type Event interface {
	EventName() string
}

type EventName string

func (e EventName) EventName() string {
	return string(e)
}

type HandleFunc func(event Event) error

type Pub struct {
	commandC chan any
	doneC    chan struct{}
	subsGC   time.Duration
}

func NewPub() Pub {
	return Pub{
		commandC: make(chan any),
		doneC:    make(chan struct{}),
		subsGC:   1 * time.Minute,
	}
}

type sub struct {
	id     int
	topics []string
	handle func(event Event) error
	doneC  chan<- struct{}
	closed bool
}

// Serve starts the publisher and blocks until context is canceled.
func (p Pub) Serve(ctx context.Context) error {
	select {
	case <-p.doneC:
		return ErrPubSubClosed
	default:
	}
	defer close(p.doneC)

	t := time.NewTicker(p.subsGC)
	defer t.Stop()

	var lastID int
	var subs []sub

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			subs = slices.DeleteFunc(subs, func(s sub) bool { return s.closed })
		case command := <-p.commandC:
			switch command := command.(type) {
			case subscribe:
				lastID++
				sub := sub{
					id:     lastID,
					topics: command.topics,
					handle: command.handle,
					doneC:  command.doneC,
					closed: false,
				}
				subs = append(subs, sub)
				command.resC <- sub.id
			case unsubscribe:
				for i := range subs {
					if subs[i].id == command.id {
						close(subs[i].doneC)
						subs[i].closed = true
					}
				}
			case Event:
				eventName := command.EventName()
				for i := range subs {
					if subs[i].closed || !slices.Contains(subs[i].topics, eventName) {
						continue
					}

					err := subs[i].handle(command)
					if err != nil {
						close(subs[i].doneC)
						subs[i].closed = true
						// TODO: handle error
					}
				}
			}
		}
	}
}

type subscribe struct {
	topics []string
	handle HandleFunc
	doneC  chan<- struct{}

	resC chan int
}

func (p Pub) Subscribe(ctx context.Context, handle HandleFunc, events ...Event) (Sub, error) {
	var topics []string
	for _, e := range events {
		topics = append(topics, e.EventName())
	}

	doneC := make(chan struct{})
	resC := make(chan int, 1)
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

type unsubscribe struct {
	id int
}

func (p Pub) Publish(ctx context.Context, event Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.doneC:
		return ErrPubSubClosed
	case p.commandC <- event:
		return nil
	}
}

type Sub struct {
	doneC <-chan struct{}
	pub   Pub
	id    int
}

// Wait blocks until the subscription is closed.
func (s *Sub) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.pub.doneC:
		return ErrPubSubClosed
	case <-s.doneC:
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
