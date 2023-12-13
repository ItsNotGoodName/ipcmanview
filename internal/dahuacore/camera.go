package dahuacore

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func NewMemCameraStore() *MemCameraStore {
	return &MemCameraStore{
		mu:   sync.Mutex{},
		data: make(map[int64]models.DahuaConn),
	}
}

type MemCameraStore struct {
	mu   sync.Mutex
	data map[int64]models.DahuaConn
}

func (s *MemCameraStore) Save(ctx context.Context, camera ...models.DahuaConn) error {
	s.mu.Lock()
	for _, camera := range camera {
		s.data[camera.ID] = camera
	}
	s.mu.Unlock()
	return nil
}

func (s *MemCameraStore) Get(ctx context.Context, id int64) (models.DahuaConn, bool, error) {
	s.mu.Lock()
	camera, found := s.data[id]
	s.mu.Unlock()
	return camera, found, nil
}

func (s *MemCameraStore) List(ctx context.Context) ([]models.DahuaConn, error) {
	s.mu.Lock()
	cameras := make([]models.DahuaConn, 0, len(s.data))
	for _, c := range s.data {
		cameras = append(cameras, c)
	}
	s.mu.Unlock()
	slices.SortFunc(cameras, func(a, b models.DahuaConn) int { return cmp.Compare(a.ID, b.ID) })
	return cameras, nil
}
