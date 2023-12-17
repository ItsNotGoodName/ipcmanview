package dahuacore

import (
	"context"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/rs/zerolog/log"
)

func cameraEqual(lhs, rhs models.DahuaConn) bool {
	return lhs.Address == rhs.Address && lhs.Username == rhs.Username && lhs.Password == rhs.Password
}

func newStoreClient(camera models.DahuaConn, lastAccessed time.Time) storeClient {
	return storeClient{
		LastAccessed: lastAccessed,
		Conn:         NewConn(camera),
	}
}

type storeClient struct {
	LastAccessed time.Time
	Conn         Conn
}

func (c storeClient) Close(ctx context.Context) {
	if err := c.Conn.RPC.Close(ctx); err != nil {
		log.Err(err).Int64("id", c.Conn.Camera.ID).Msg("Failed to close RPC connection")
	}
}

type Store struct {
	clientsMu sync.Mutex
	clients   map[int64]storeClient
}

func NewStore() *Store {
	return &Store{
		clientsMu: sync.Mutex{},
		clients:   make(map[int64]storeClient),
	}
}

func (s *Store) Serve(ctx context.Context) error {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			wg := sync.WaitGroup{}

			s.clientsMu.Lock()
			for _, client := range s.clients {
				wg.Add(1)
				go func(client storeClient) {
					client.Close(context.Background())
					wg.Done()
				}(client)
			}
			s.clientsMu.Unlock()

			wg.Wait()

			return ctx.Err()
		case <-t.C:
			var clients []storeClient

			now := time.Now()

			s.clientsMu.Lock()
			for id, client := range s.clients {
				if now.Sub(client.LastAccessed) > 5*time.Minute {
					delete(s.clients, id)
					clients = append(clients, client)
				}
			}
			s.clientsMu.Unlock()

			for _, client := range clients {
				client.Close(ctx)
			}
		}
	}
}

func (s *Store) getOrCreateCamera(ctx context.Context, camera models.DahuaConn) Conn {
	client, ok := s.clients[camera.ID]
	if !ok {
		// Not found

		client = newStoreClient(camera, time.Now())
		s.clients[camera.ID] = client
	} else if !cameraEqual(client.Conn.Camera, camera) {
		// Found but not equal

		// Closing camera connection should not block that store
		go client.Close(ctx)

		client = newStoreClient(camera, time.Now())
		s.clients[camera.ID] = client
	} else {
		// Found

		client.LastAccessed = time.Now()
		s.clients[camera.ID] = client
	}

	return client.Conn
}

func (s *Store) ConnList(ctx context.Context, cameras []models.DahuaConn) []Conn {
	clients := make([]Conn, 0, len(cameras))

	s.clientsMu.Lock()
	for _, camera := range cameras {
		clients = append(clients, s.getOrCreateCamera(ctx, camera))
	}
	s.clientsMu.Unlock()

	return clients
}

func (s *Store) Conn(ctx context.Context, camera models.DahuaConn) Conn {
	s.clientsMu.Lock()
	client := s.getOrCreateCamera(ctx, camera)
	s.clientsMu.Unlock()

	return client
}

func (s *Store) ConnDelete(ctx context.Context, id int64) {
	s.clientsMu.Lock()
	client, found := s.clients[id]
	if found {
		delete(s.clients, id)
	}
	s.clientsMu.Unlock()

	if found {
		client.Close(ctx)
	}
}

func (store *Store) Register(bus *core.Bus) {
	bus.OnEventDahuaCameraDeleted(func(ctx context.Context, evt models.EventDahuaCameraDeleted) error {
		store.ConnDelete(ctx, evt.CameraID)
		return nil
	})
}
