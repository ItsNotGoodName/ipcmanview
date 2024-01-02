package dahua

import (
	"context"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/rs/zerolog/log"
)

func newStoreClient(conn models.DahuaConn, lastAccessed time.Time) storeClient {
	return storeClient{
		LastAccessed: lastAccessed,
		Client:       NewClient(conn),
	}
}

type storeClient struct {
	LastAccessed time.Time
	Client       Client
}

func (c storeClient) Close(ctx context.Context) {
	if err := c.Client.RPC.Close(ctx); err != nil {
		log.Err(err).Int64("id", c.Client.Conn.ID).Msg("Failed to close RPC connection")
	}
}

// Store creates and maintains clients to devices.
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

func (*Store) String() string {
	return "dahua.Store"
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
				if now.Sub(client.LastAccessed) > 5*time.Minute && now.Sub(client.Client.RPC.State(ctx).LastRPC) > 5*time.Minute {
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

func (s *Store) getOrCreateClient(ctx context.Context, conn models.DahuaConn) Client {
	client, ok := s.clients[conn.ID]
	if !ok {
		// Not found

		client = newStoreClient(conn, time.Now())
		s.clients[conn.ID] = client
	} else if !ConnEqual(client.Client.Conn, conn) {
		// Found but not equal

		// Closing device connection should not block that store
		go client.Close(ctx)

		client = newStoreClient(conn, time.Now())
		s.clients[conn.ID] = client
	} else {
		// Found

		client.LastAccessed = time.Now()
		s.clients[conn.ID] = client
	}

	return client.Client
}

func (s *Store) ClientList(ctx context.Context, conns []models.DahuaConn) []Client {
	clients := make([]Client, 0, len(conns))

	s.clientsMu.Lock()
	for _, conn := range conns {
		clients = append(clients, s.getOrCreateClient(ctx, conn))
	}
	s.clientsMu.Unlock()

	return clients
}

func (s *Store) Client(ctx context.Context, conn models.DahuaConn) Client {
	s.clientsMu.Lock()
	client := s.getOrCreateClient(ctx, conn)
	s.clientsMu.Unlock()

	return client
}

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

func (s *Store) Register(bus *core.Bus) *Store {
	bus.OnEventDahuaDeviceDeleted(func(ctx context.Context, evt models.EventDahuaDeviceDeleted) error {
		s.ClientDelete(ctx, evt.DeviceID)
		return nil
	})
	return s
}
