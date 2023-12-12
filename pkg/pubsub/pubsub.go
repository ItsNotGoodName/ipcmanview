// pubsub is a simple in-memory event pub sub.
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
	errC   chan<- error
	closed bool
}

type subscribe struct {
	topics []string
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error

	resC chan int
}

type unsubscribe struct {
	id int
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
			for i := range subs {
				if subs[i].closed {
					subs = slices.DeleteFunc(subs, func(s sub) bool { return s.closed })
					break
				}
			}
		case command := <-p.commandC:
			switch command := command.(type) {
			case subscribe:
				lastID++
				sub := sub{
					id:     lastID,
					topics: command.topics,
					handle: command.handle,
					doneC:  command.doneC,
					errC:   command.errC,
					closed: false,
				}
				subs = append(subs, sub)
				command.resC <- sub.id
			case unsubscribe:
				for i := range subs {
					if subs[i].id == command.id {
						subs[i].errC <- nil
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
						subs[i].errC <- err
						close(subs[i].doneC)
						subs[i].closed = true
					}
				}
			}
		}
	}
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
