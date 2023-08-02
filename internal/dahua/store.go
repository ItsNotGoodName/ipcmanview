package dahua

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
)

type StoreActor ActorHandle

func (c StoreActor) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	return ActorHandle(c).RPC(ctx)
}

type Store struct {
	mu     sync.Mutex
	actors []StoreActor
}

func NewStore() *Store {
	return &Store{
		actors: []StoreActor{},
	}
}

func (s *Store) GetOrCreate(dbCtx db.Context, id int64) (StoreActor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cam, err := db.DahuaCameraGet(dbCtx, id)
	if err != nil {
		return StoreActor{}, err
	}

	for i := range s.actors {
		if s.actors[i].Camera.ID != id {
			continue
		}

		if !s.actors[i].Camera.Equal(cam) {
			ActorHandle(s.actors[i]).Close(dbCtx.Context)
			s.actors[i] = StoreActor(NewActorHandle(cam))
		}

		return s.actors[i], nil
	}

	actor := StoreActor(NewActorHandle(cam))
	s.actors = append(s.actors, actor)

	return actor, nil
}

func (s *Store) Delete(ctx context.Context, id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	actors := []StoreActor{}
	for i := range s.actors {
		if s.actors[i].Camera.ID != id {
			actors = append(actors, StoreActor(s.actors[i]))
			continue
		}

		ActorHandle(s.actors[i]).Close(ctx)
	}

	s.actors = actors
}
