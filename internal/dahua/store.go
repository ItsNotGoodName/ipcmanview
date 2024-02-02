package dahua

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func newStoreClient(conn models.Conn) storeClient {
	return storeClient{
		Client: NewClient(conn),
	}
}

type storeClient struct {
	Client Client
}

func (c storeClient) Close(ctx context.Context) {
	if err := c.Client.Close(ctx); err != nil {
		log.Err(err).Int64("id", c.Client.Conn.ID).Msg("Failed to close RPC connection")
	}
}

func NewStore(db repo.DB) *Store {
	return &Store{
		ServiceContext: sutureext.NewServiceContext("dahua.Store"),
		db:             db,
		clientsMu:      sync.Mutex{},
		clients:        make(map[int64]storeClient),
	}
}

// Store holds dahua clients.
type Store struct {
	sutureext.ServiceContext
	db        repo.DB
	clientsMu sync.Mutex
	clients   map[int64]storeClient
}

// Close closes all clients.
func (s *Store) Close() {
	wg := sync.WaitGroup{}

	for _, client := range s.clients {
		wg.Add(1)
		go func(client storeClient) {
			client.Close(s.Context())
			wg.Done()
		}(client)
	}

	wg.Wait()
}

func (s *Store) getOrCreateClient(ctx context.Context, conn models.Conn) Client {
	client, ok := s.clients[conn.ID]
	if !ok {
		// Not found

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	} else if !client.Client.Conn.EQ(conn) {
		// Found but not equal

		// Closing device connection should not block the store
		go client.Close(s.Context())

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	} else {
		// Found

		s.clients[conn.ID] = client
	}

	return client.Client
}

func (s *Store) GetClient(ctx context.Context, id int64) (Client, error) {
	s.clientsMu.Lock()
	conn, err := s.db.DahuaGetConn(ctx, id)
	if err != nil {
		s.clientsMu.Unlock()
		return Client{}, err
	}

	client := s.getOrCreateClient(ctx, conn)
	s.clientsMu.Unlock()

	return client, nil
}

func (s *Store) ListClient(ctx context.Context, ids []int64) ([]Client, error) {
	s.clientsMu.Lock()
	conns, err := s.db.DahuaListConn(ctx, ids)
	if err != nil {
		s.clientsMu.Unlock()
		return nil, err
	}

	var clients []Client
	for _, conn := range conns {
		clients = append(clients, s.getOrCreateClient(ctx, conn))
	}
	s.clientsMu.Unlock()

	return clients, nil
}

func (s *Store) Register(bus *event.Bus) *Store {
	bus.OnDahuaDeviceUpdated(func(ctx context.Context, evt event.DahuaDeviceUpdated) error {
		s.clientsMu.Lock()
		s.getOrCreateClient(ctx, evt.Conn)
		s.clientsMu.Unlock()
		return nil
	})
	bus.OnDahuaDeviceDeleted(func(ctx context.Context, evt event.DahuaDeviceDeleted) error {
		s.clientsMu.Lock()
		client, found := s.clients[evt.DeviceID]
		if found {
			delete(s.clients, evt.DeviceID)
		}
		s.clientsMu.Unlock()

		if found {
			go client.Close(s.Context())
		}
		return nil
	})
	return s
}
