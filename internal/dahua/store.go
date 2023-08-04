package dahua

import (
	"context"
	"errors"
	"sync"

	"github.com/ItsNotGoodName/ipcmango/internal/event"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type StoreActorHandle ActorHandle

func (c StoreActorHandle) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	return ActorHandle(c).RPC(ctx)
}

type Store struct {
	mu     sync.Mutex
	actors []StoreActorHandle
}

func NewStore() *Store {
	return &Store{
		actors: []StoreActorHandle{},
	}
}

func (s *Store) GetOrCreate(ctx context.Context, db qes.Querier, id int64) (StoreActorHandle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cam, err := DB.CameraGet(ctx, db, id)
	if err != nil {
		return StoreActorHandle{}, err
	}

	for i := range s.actors {
		if s.actors[i].cam.ID != id {
			continue
		}

		if !s.actors[i].cam.Equal(cam) {
			ActorHandle(s.actors[i]).Close(ctx)
			s.actors[i] = StoreActorHandle(NewActorHandle(cam))
		}

		return s.actors[i], nil
	}

	actor := StoreActorHandle(NewActorHandle(cam))
	s.actors = append(s.actors, actor)

	return actor, nil
}

func (s *Store) Delete(ctx context.Context, id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	actors := []StoreActorHandle{}
	for i := range s.actors {
		if s.actors[i].cam.ID != id {
			actors = append(actors, StoreActorHandle(s.actors[i]))
			continue
		}

		ActorHandle(s.actors[i]).Close(ctx)
	}

	s.actors = actors
}

// TODO: I don't like how this is named
func StoreConnectBus(bus *event.Bus, store *Store, pool *pgxpool.Pool) {
	bus.DahuaCameraUpdated = append(bus.DahuaCameraUpdated, func(ctx context.Context, evt event.DahuaCameraUpdated) {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			return
		}
		defer conn.Release()

		for _, v := range evt.IDS {
			_, err := store.GetOrCreate(ctx, conn, v)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Err(err).Msg("Failed to update dahua store's camera")
			}
		}
	})
	bus.DahuaCameraDeleted = append(bus.DahuaCameraDeleted, func(ctx context.Context, evt event.DahuaCameraDeleted) {
		for _, v := range evt.IDS {
			store.Delete(ctx, v)
		}
	})
}
