package dahua

import (
	"context"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/rs/zerolog/log"
)

func newStoreClient(conn Conn) storeClient {
	return storeClient{
		Client: NewClient(conn),
	}
}

type storeClient struct {
	Client
}

func (c storeClient) logError(err error) {
	if err != nil {
		log.Err(err).Int64("id", c.Client.Conn.ID).Msg("Failed to close Client connection")
	}
}

func NewStore() *Store {
	return &Store{
		shutdownTimeout: 3 * time.Second,
		clientsMu:       sync.Mutex{},
		clients:         make(map[int64]storeClient),
	}
}

// Store holds clients.
type Store struct {
	shutdownTimeout time.Duration

	clientsMu sync.Mutex
	clients   map[int64]storeClient
}

func (*Store) String() string {
	return "dahua.Store"
}

// Close closes all clients.
func (s *Store) Close() {
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	for _, client := range s.clients {
		wg.Add(1)
		go func(client storeClient) {
			defer wg.Done()
			client.logError(client.Client.Close(ctx))
		}(client)
	}

	wg.Wait()
}

func (s *Store) getOrCreateClient(ctx context.Context, conn Conn) Client {
	client, ok := s.clients[conn.ID]
	if !ok {
		// Not found

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	} else if !client.Client.Conn.EQ(conn) {
		// Found but not equal

		client.logError(client.Client.CloseNoWait(ctx))

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	}

	return client.Client
}

func (s *Store) GetClient(ctx context.Context, deviceID int64) (Client, error) {
	s.clientsMu.Lock()
	conn, err := GetConn(ctx, deviceID)
	if err != nil {
		s.clientsMu.Unlock()
		return Client{}, err
	}

	client := s.getOrCreateClient(ctx, conn)
	s.clientsMu.Unlock()

	return client, nil
}

func (s *Store) ListClient(ctx context.Context) ([]Client, error) {
	s.clientsMu.Lock()
	conns, err := ListConn(ctx)
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

func (s *Store) deleteClient(ctx context.Context, deviceID int64) {
	s.clientsMu.Lock()
	client, found := s.clients[deviceID]
	if found {
		delete(s.clients, deviceID)
	}
	s.clientsMu.Unlock()

	if found {
		client.logError(client.Client.Close(ctx))
	}
}

func (s *Store) Register(hub *bus.Hub) *Store {
	upsert := func(ctx context.Context, deviceID int64) error {
		if _, err := s.GetClient(ctx, deviceID); err != nil {
			if core.IsNotFound(err) {
				s.deleteClient(ctx, deviceID)
				return nil
			}
			return err
		}
		return nil
	}

	hub.OnDahuaDeviceCreated(s.String(), func(ctx context.Context, event bus.DahuaDeviceCreated) error {
		return upsert(ctx, event.DeviceID)
	})
	hub.OnDahuaDeviceUpdated(s.String(), func(ctx context.Context, event bus.DahuaDeviceUpdated) error {
		return upsert(ctx, event.DeviceID)
	})
	hub.OnDahuaDeviceDeleted(s.String(), func(ctx context.Context, event bus.DahuaDeviceDeleted) error {
		s.deleteClient(ctx, event.DeviceID)
		return nil
	})

	return s
}
