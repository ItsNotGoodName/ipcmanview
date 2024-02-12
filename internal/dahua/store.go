package dahua

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
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

func NewStore(db sqlite.DB) *Store {
	return &Store{
		ServiceContext: sutureext.NewServiceContext("dahua.Store"),
		db:             db,
		clientsMu:      sync.Mutex{},
		clients:        make(map[int64]storeClient),
	}
}

// Store holds clients.
type Store struct {
	sutureext.ServiceContext
	db        sqlite.DB
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

func (s *Store) getOrCreateClient(ctx context.Context, conn Conn) Client {
	client, ok := s.clients[conn.ID]
	if !ok {
		// Not found

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	} else if !client.Client.Conn.EQ(conn) {
		// Found but not equal

		// Closing client should not block the store
		go client.Close(s.Context())

		client = newStoreClient(conn)
		s.clients[conn.ID] = client
	}

	return client.Client
}

func (s *Store) GetClient(ctx context.Context, deviceID int64) (Client, error) {
	s.clientsMu.Lock()
	conn, err := GetConn(ctx, s.db, deviceID)
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
	conns, err := ListConn(ctx, s.db)
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

func (s *Store) deleteClient(ctx context.Context, deviceID int64) error {
	s.clientsMu.Lock()
	client, found := s.clients[deviceID]
	if found {
		delete(s.clients, deviceID)
	}
	s.clientsMu.Unlock()

	if found {
		client.Close(s.Context())
	}
	return nil
}

func (s *Store) Register(bus *event.Bus) *Store {
	bus.OnEvent(func(ctx context.Context, evt event.Event) error {
		switch evt.Event.Action {
		case event.ActionDahuaDeviceCreated, event.ActionDahuaDeviceUpdated:
			deviceID := event.DataAsInt64(evt.Event)

			if _, err := s.GetClient(ctx, deviceID); err != nil {
				if repo.IsNotFound(err) {
					return s.deleteClient(ctx, deviceID)
				}
				return err
			}
		case event.ActionDahuaDeviceDeleted:
			deviceID := event.DataAsInt64(evt.Event)

			return s.deleteClient(ctx, deviceID)
		}
		return nil
	})
	return s
}
