package dahuacore

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func NewMemDeviceStore() *MemDeviceStore {
	return &MemDeviceStore{
		mu:   sync.Mutex{},
		data: make(map[int64]models.DahuaConn),
	}
}

type MemDeviceStore struct {
	mu   sync.Mutex
	data map[int64]models.DahuaConn
}

func (s *MemDeviceStore) Save(ctx context.Context, device ...models.DahuaConn) error {
	s.mu.Lock()
	for _, device := range device {
		s.data[device.ID] = device
	}
	s.mu.Unlock()
	return nil
}

func (s *MemDeviceStore) Get(ctx context.Context, id int64) (models.DahuaConn, bool, error) {
	s.mu.Lock()
	device, found := s.data[id]
	s.mu.Unlock()
	return device, found, nil
}

func (s *MemDeviceStore) List(ctx context.Context) ([]models.DahuaConn, error) {
	s.mu.Lock()
	devices := make([]models.DahuaConn, 0, len(s.data))
	for _, c := range s.data {
		devices = append(devices, c)
	}
	s.mu.Unlock()
	slices.SortFunc(devices, func(a, b models.DahuaConn) int { return cmp.Compare(a.ID, b.ID) })
	return devices, nil
}
