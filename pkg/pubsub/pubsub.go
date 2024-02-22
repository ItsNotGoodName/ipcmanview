// Package pubsub is a simple in-memory pub sub.
package pubsub

import (
	"context"
	"slices"
	"sync"
)

type Event interface {
	EventTopic() string
}

type EventTopic string

func (e EventTopic) EventTopic() string {
	return string(e)
}

type HandleFunc func(ctx context.Context, evt Event) error

type MiddlewareFunc func(next HandleFunc) HandleFunc

type StateSubscriber struct {
	ID     int
	Topics []string
}

type State struct {
	SubscriberCount int
	Subscribers     []StateSubscriber
}

type sub struct {
	topics []string
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error
}

type Pub struct {
	middleware []MiddlewareFunc

	mu     sync.Mutex
	lastID int
	subs   map[int]sub
}

func NewPub(middleware ...MiddlewareFunc) *Pub {
	return &Pub{
		middleware: middleware,
		mu:         sync.Mutex{},
		lastID:     0,
		subs:       make(map[int]sub),
	}
}

func (p *Pub) Publish(ctx context.Context, event Event) error {
	p.mu.Lock()
	for i := range p.subs {
		if !slices.Contains(p.subs[i].topics, event.EventTopic()) {
			continue
		}

		if err := p.subs[i].handle(ctx, event); err != nil {
			p.subs[i].errC <- err
			close(p.subs[i].doneC)
			delete(p.subs, i)
		}
	}
	p.mu.Unlock()

	return nil
}

type subscribeParams struct {
	topics []string
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error
}

func (p *Pub) subscribe(arg subscribeParams) int {
	p.mu.Lock()
	p.lastID++
	id := p.lastID
	p.subs[id] = sub{
		topics: arg.topics,
		handle: arg.handle,
		doneC:  arg.doneC,
		errC:   arg.errC,
	}
	p.mu.Unlock()

	return id
}

func (p *Pub) unsubscribe(id int, err error) {
	p.mu.Lock()
	sub, ok := p.subs[id]
	if !ok {
		p.mu.Unlock()
		return
	}

	sub.errC <- err
	close(sub.doneC)
	delete(p.subs, id)
	p.mu.Unlock()
}

func (p *Pub) State() (State, error) {
	p.mu.Lock()
	ss := make([]StateSubscriber, 0, len(p.subs))
	ssc := 0
	for i := range p.subs {
		ssc++
		ss = append(ss, StateSubscriber{
			ID:     i,
			Topics: p.subs[i].topics,
		})
	}
	s := State{
		SubscriberCount: ssc,
		Subscribers:     ss,
	}
	p.mu.Unlock()

	return s, nil
}

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

func (s Sub) Close() {
	s.pub.unsubscribe(s.id, nil)
}

type SubscribeBuilder struct {
	pub        *Pub
	topics     []string
	middleware []MiddlewareFunc
}

func (p *Pub) Subscribe(events ...Event) SubscribeBuilder {
	// Get topics
	var topics []string
	for _, e := range events {
		topics = append(topics, e.EventTopic())
	}

	// Copy middleware
	middleware := make([]MiddlewareFunc, 0, len(p.middleware))
	copy(middleware, p.middleware)

	return SubscribeBuilder{
		pub:        p,
		topics:     topics,
		middleware: middleware,
	}
}

// Middleware add a middleware between the handler and publisher.
func (b SubscribeBuilder) Middleware(fn MiddlewareFunc) SubscribeBuilder {
	b.middleware = append(b.middleware, fn)
	return b
}

// Function creates a subscription with a function.
func (b SubscribeBuilder) Function(fn HandleFunc) (Sub, error) {
	handle := fn

	// Attach middleware
	for _, mw := range b.middleware {
		handle = mw(handle)
	}

	// Subscribe
	doneC := make(chan struct{})
	errC := make(chan error, 1)
	id := b.pub.subscribe(subscribeParams{
		topics: b.topics,
		handle: handle,
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
	evtC := make(chan Event, size)

	// Subscribe
	sub, err := b.Function(func(pubCtx context.Context, evt Event) error {
		select {
		case <-pubCtx.Done():
			return pubCtx.Err()
		case <-ctx.Done():
			return ctx.Err()
		case evtC <- evt:
			return nil
		}
	})
	if err != nil {
		return Sub{}, nil, err
	}

	// Close event channel
	go func() {
		<-sub.doneC
		close(evtC)
	}()

	return sub, evtC, nil
}
