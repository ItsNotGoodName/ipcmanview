package dahua

import (
	"sync"

	"github.com/ItsNotGoodName/ipcmango/internal/db"
)

type Store struct {
	mu     sync.Mutex
	actors []CameraActor
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Get(context db.Context, id int64) (CameraActor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, actor := range s.actors {
		if actor.ID == id {
			return actor, nil
		}
	}

	cam, err := db.DahuaCameraGet(context, id)
	if err != nil {
		return CameraActor{}, err
	}

	actor := NewCameraActor(cam)
	s.actors = append(s.actors, actor)

	return actor, nil
}
