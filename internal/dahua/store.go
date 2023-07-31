package dahua

import (
	"context"
	"errors"
	"sync"

	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
)

type StoreActor Actor

func (c StoreActor) RPC(ctx context.Context) (dahua.RequestBuilder, error) {
	return Actor(c).RPC(ctx)
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

func (s *Store) GetOrCreate(context db.Context, id int64) (StoreActor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cam, err := db.DahuaCameraGet(context, id)
	if err != nil {
		return StoreActor{}, err
	}

	for _, actor := range s.actors {
		if actor.ID != id {
			continue
		}

		err := Actor(actor).Update(context.Context, cam)
		if err != nil {
			if errors.Is(err, ErrActorClosed) {
				break
			}

			return StoreActor{}, err
		}

		return actor, nil
	}

	actor := StoreActor(StartActor(cam))
	s.actors = append(s.actors, actor)
	return actor, nil
}

func (s *Store) Delete(ctx context.Context, id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	actors := []StoreActor{}
	for i := range s.actors {
		if s.actors[i].ID != id {
			actors = append(actors, StoreActor(s.actors[i]))
			continue
		}

		Actor(s.actors[i]).Close(ctx)
	}

	s.actors = actors
}
