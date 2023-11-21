package pubsub

import (
	"context"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

type pubSub struct {
	Ch  chan<- models.EventDahuaCameraEvent
	IDs []string
}

func NewPub(dahuaBus DahuaBus) *Pub {
	return register(dahuaBus, &Pub{
		mu:                    sync.Mutex{},
		id:                    0,
		dahuaEventSubscribers: map[int]pubSub{},
	})
}

type Pub struct {
	mu                    sync.Mutex
	id                    int
	dahuaEventSubscribers map[int]pubSub
}

func (p *Pub) unsubscribeDahuaEvents(id int) {
	sub, found := p.dahuaEventSubscribers[id]
	if found {
		close(sub.Ch)
		delete(p.dahuaEventSubscribers, id)
	}
}

func (p *Pub) SubscribeDahuaEvents(ctx context.Context, ids []string) (<-chan models.EventDahuaCameraEvent, error) {
	ch := make(chan models.EventDahuaCameraEvent, 100)
	slices.Sort(ids)

	p.mu.Lock()
	p.id += 1
	id := p.id
	p.dahuaEventSubscribers[id] = pubSub{
		Ch:  ch,
		IDs: ids,
	}
	p.mu.Unlock()

	go func() {
		<-ctx.Done()
		p.mu.Lock()
		p.unsubscribeDahuaEvents(id)
		p.mu.Unlock()
	}()

	return ch, nil
}

type DahuaBus interface {
	OnCameraEvent(h func(ctx context.Context, evt models.EventDahuaCameraEvent) error)
}

func register(dahuaBus DahuaBus, p *Pub) *Pub {
	dahuaBus.OnCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		p.mu.Lock()
		defer p.mu.Unlock()

		var deadSubIDs []int
		for id, sub := range p.dahuaEventSubscribers {
			if len(sub.IDs) != 0 {
				_, found := slices.BinarySearch(sub.IDs, evt.Event.CameraID)
				if !found {
					continue
				}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case sub.Ch <- evt:
			default:
				deadSubIDs = append(deadSubIDs, id)
			}
		}
		for _, id := range deadSubIDs {
			p.unsubscribeDahuaEvents(id)
		}

		return nil
	})
	return p
}
