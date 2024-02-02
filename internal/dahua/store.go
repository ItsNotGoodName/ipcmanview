package dahua

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog/log"
)

func newStoreClient(conn Conn) storeClient {
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

func NewStore() *Store {
	return &Store{
		ServiceContext: sutureext.NewServiceContext("dahua.Store"),
		clientsMu:      sync.Mutex{},
		clients:        make(map[int64]storeClient),
	}
}

// Store deduplicates device clients.
type Store struct {
	sutureext.ServiceContext
	clientsMu sync.Mutex
	clients   map[int64]storeClient
}

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

func (s *Store) getOrCreateClient(ctx context.Context, conn Conn) Client {
	client, ok := s.clients[conn.ID]
	if !ok {
		// Not found

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	} else if !client.Client.Conn.EQ(conn) && conn.UpdatedAt.After(client.Client.Conn.UpdatedAt) {
		// Found but not equal and newer

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

func (s *Store) ClientList(ctx context.Context, conns []Conn) []Client {
	clients := make([]Client, 0, len(conns))

	s.clientsMu.Lock()
	for _, conn := range conns {
		clients = append(clients, s.getOrCreateClient(ctx, conn))
	}
	s.clientsMu.Unlock()

	return clients
}

func (s *Store) Client(ctx context.Context, conn Conn) Client {
	s.clientsMu.Lock()
	client := s.getOrCreateClient(ctx, conn)
	s.clientsMu.Unlock()

	return client
}

// FIXME: deleted clients are recreated when an old connection is passed to Store.Client or Store.ClientList
func (s *Store) ClientDelete(ctx context.Context, id int64) {
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

func (s *Store) Register(bus *event.Bus) *Store {
	bus.OnDahuaDeviceDeleted(func(ctx context.Context, evt event.DahuaDeviceDeleted) error {
		s.ClientDelete(ctx, evt.DeviceID)
		return nil
	})
	return s
}
