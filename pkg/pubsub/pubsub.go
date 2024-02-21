// Package pubsub is a simple in-memory event pub sub.
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

type StateSubscriber struct {
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
	mu     sync.Mutex
	lastID int
	subs   map[int]sub
}

func NewPub() *Pub {
	return &Pub{
		mu:     sync.Mutex{},
		lastID: 0,
		subs:   make(map[int]sub),
	}
}

func (p *Pub) Publish(ctx context.Context, event Event) error {
	p.mu.Lock()
	for i := range p.subs {
		if !slices.Contains(p.subs[i].topics, event.EventTopic()) {
			continue
		}

		if err := p.subs[i].handle(ctx, event); err != nil {
			p.subs[i].errC <- nil
			close(p.subs[i].doneC)
			delete(p.subs, i)
		}
	}
	p.mu.Unlock()

	return nil
}

type subscribe struct {
	topics []string
	handle HandleFunc
	doneC  chan<- struct{}
	errC   chan<- error

	resC chan<- int
}

func (p *Pub) subscribe(arg subscribe) int {
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

func (p *Pub) unsubscribe(id int) {
	p.mu.Lock()
	for i := range p.subs {
		if i == id {
			p.subs[i].errC <- nil
			close(p.subs[i].doneC)
			delete(p.subs, i)
		}
	}
	p.mu.Unlock()
}

func (p *Pub) State(ctx context.Context) (State, error) {
	p.mu.Lock()
	ss := make([]StateSubscriber, 0, len(p.subs))
	ssc := 0
	for i := range p.subs {
		ssc++
		ss = append(ss, StateSubscriber{
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
