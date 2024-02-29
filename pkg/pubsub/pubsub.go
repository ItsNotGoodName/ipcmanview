// Package pubsub is a simple in-memory pub sub.
package pubsub

import (
	"context"
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
	ID int
}

type State struct {
	SubscriberCount int
	Subscribers     []StateSubscriber
}

type sub struct {
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

func (p *Pub) Broadcast(ctx context.Context, event Event) error {
	p.mu.Lock()
	for i := range p.subs {
		if err := p.subs[i].handle(ctx, event); err != nil {
			// Remove subscriber
			p.subs[i].errC <- err
			close(p.subs[i].doneC)
			delete(p.subs, i)
		}
	}
	p.mu.Unlock()

	return nil
}

func (p *Pub) unsubscribe(id int, err error) {
	p.mu.Lock()
	sub, ok := p.subs[id]
	if !ok {
		p.mu.Unlock()
		return
	}

	// Remove subscriber
	sub.errC <- err
	close(sub.doneC)
	delete(p.subs, id)
	p.mu.Unlock()
}

type subscribeParams struct {
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error
}

func (p *Pub) subscribe(arg subscribeParams) int {
	p.mu.Lock()
	p.lastID++
	id := p.lastID
	p.subs[id] = sub{
		handle: arg.handle,
		doneC:  arg.doneC,
		errC:   arg.errC,
	}
	p.mu.Unlock()

	return id
}

func (p *Pub) State() (State, error) {
	p.mu.Lock()
	ss := make([]StateSubscriber, 0, len(p.subs))
	ssc := 0
	for i := range p.subs {
		ssc++
		ss = append(ss, StateSubscriber{
			ID: i,
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
	middleware []MiddlewareFunc
}

func (p *Pub) Subscribe() SubscribeBuilder {
	// Copy middleware
	middleware := make([]MiddlewareFunc, 0, len(p.middleware))
	copy(middleware, p.middleware)

	return SubscribeBuilder{
		pub:        p,
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
