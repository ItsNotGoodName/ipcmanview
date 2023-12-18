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
	EventTopic() string
}

type EventTopic string

func (e EventTopic) EventTopic() string {
	return string(e)
}

type HandleFunc func(ctx context.Context, event Event) error

type StateSubscriber struct {
	Topics []string
}

type State struct {
	SubscriberCount int
	Subscribers     []StateSubscriber
	LastGC          time.Time
}

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
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error
	closed bool
}

type subscribe struct {
	topics []string
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error

	resC chan<- int
}

type unsubscribe struct {
	id int
}

type state struct {
	resC chan<- State
}

func (Pub) String() string {
	return "pubsub.Pub"
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
	defer func() {
		for i := range subs {
			if !subs[i].closed {
				subs[i].errC <- ErrPubSubClosed
				close(subs[i].doneC)
				subs[i].closed = true
			}
		}
	}()

	var lastGC time.Time

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
			lastGC = time.Now()
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
			case state:
				stateSubs := make([]StateSubscriber, 0, len(subs))
				stateSubCount := 0
				for i := range subs {
					if subs[i].closed {
						continue
					}
					stateSubCount++
					stateSubs = append(stateSubs, StateSubscriber{
						Topics: subs[i].topics,
					})
				}
				s := State{
					SubscriberCount: stateSubCount,
					Subscribers:     stateSubs,
					LastGC:          lastGC,
				}
				command.resC <- s
			case Event:
				eventTopic := command.EventTopic()
				for i := range subs {
					if subs[i].closed || !slices.Contains(subs[i].topics, eventTopic) {
						continue
					}

					err := subs[i].handle(ctx, command)
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

func (p Pub) State(ctx context.Context) (State, error) {
	resC := make(chan State, 1)

	select {
	case <-ctx.Done():
		return State{}, ctx.Err()
	case <-p.doneC:
		return State{}, ErrPubSubClosed
	case p.commandC <- state{resC: resC}:
	}

	select {
	case <-ctx.Done():
		return State{}, ctx.Err()
	case <-p.doneC:
		return State{}, ErrPubSubClosed
	case s := <-resC:
		return s, nil
	}
}
