package pubsub

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

type Pub struct {
	*mqtt.Server
	id int32
}

func (p *Pub) NextID() int {
	return int(atomic.AddInt32(&p.id, 1))
}

func NewPub(bind bool, address string) (*Pub, error) {
	mqttServer := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	if bind {
		if err := mqttServer.AddListener(listeners.NewTCP("t1", address, nil)); err != nil {
			return nil, err
		}
	}
	if err := mqttServer.AddHook(new(auth.Hook), &auth.Options{
		Ledger: &auth.Ledger{
			Auth: auth.AuthRules{
				{Username: "", Password: "", Allow: true},
			},
			ACL: auth.ACLRules{
				{
					Filters: auth.Filters{
						"#": auth.ReadOnly,
					},
				},
			},
		},
	}); err != nil {
		return nil, err
	}

	return &Pub{mqttServer, 0}, nil
}

func (p *Pub) Serve(ctx context.Context) error {
	err := p.Server.Serve()
	if err != nil {
		return errors.Join(suture.ErrTerminateSupervisorTree, err)
	}

	<-ctx.Done()

	return p.Server.Close()
}

func (pub *Pub) Register(dahuaBus *dahua.Bus) {
	dahuaBus.OnCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		return SendDahuaEvent(pub, evt.Event)
	})
}

type Thing struct {
	Server *mqtt.Server
}

// SubscribeDahuaEvents implements api.PubSub.
func (Thing) SubscribeDahuaEvents(ctx context.Context, cameraIDs []int64) (<-chan models.EventDahuaCameraEvent, error) {
	panic("unimplemented")
}

type Sub struct {
	pub       *Pub
	fn        func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) error
	closeOnce sync.Once
	errC      chan error
	closedC   chan struct{}

	subsMu sync.Mutex
	subs   []sub
}

type sub struct {
	id     int
	filter string
}

func (pub *Pub) NewSub(fn func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) error) *Sub {
	return &Sub{
		pub:       pub,
		fn:        fn,
		closeOnce: sync.Once{},
		errC:      make(chan error, 1),
		closedC:   make(chan struct{}),
		subs:      []sub{},
	}
}

func (s *Sub) close(err error) {
	s.closeOnce.Do(func() {
		s.subsMu.Lock()
		for _, sub := range s.subs {
			err := s.pub.Unsubscribe(sub.filter, sub.id)
			if err != nil {
				log.Err(err).Str("topic", sub.filter).Msg("Failed to unsubscribe")
			}
		}
		close(s.closedC)
		s.errC <- err
		s.subsMu.Unlock()
	})
}

func (s *Sub) Close() {
	s.close(nil)
}

func (s *Sub) Subscribe(filter string) error {
	s.subsMu.Lock()
	defer s.subsMu.Unlock()

	select {
	case <-s.closedC:
		return errors.New("subscription closed")
	default:
	}

	id := s.pub.NextID()

	err := s.pub.Subscribe(filter, id, s.handle)
	if err != nil {
		return err
	}

	s.subs = append(s.subs, sub{
		id:     id,
		filter: filter,
	})

	return nil
}

func (s *Sub) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-s.errC:
		s.errC <- err
		return err
	}
}

func (s *Sub) handle(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	select {
	case <-s.closedC:
		return
	default:
	}
	if err := s.fn(cl, sub, pk); err != nil {
		s.close(err)
	}
}
