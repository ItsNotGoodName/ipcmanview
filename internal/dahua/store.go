package dahua

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/jmoiron/sqlx"
)

const StoreClientCloseTimeout = 3 * time.Second

func checkStoreClientCloseError(conn Conn, err error) {
	if err != nil {
		slog.Error("Failed to close Client connection", slog.String("name", conn.Name))
	}
}

func NewStoreClient(conn Conn) StoreClient {
	urL, _ := url.Parse("http://" + conn.IP)

	httpClient := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 5*time.Second)
			},
		},
	}

	clientRPC := dahuarpc.NewClient(
		"",
		conn.Username,
		conn.Password,
		dahuarpc.WithHTTPClient(&httpClient),
		dahuarpc.WithURL(dahuarpc.URL(urL), dahuarpc.LoginURL(urL)),
	)
	clientPTZ := ptz.NewClient(clientRPC)
	clientCGI := dahuacgi.NewClient(
		"",
		conn.Username,
		conn.Password,
		dahuacgi.WithURL(dahuacgi.URL(urL)),
	)
	clientFile := dahuarpc.NewFileClient(&httpClient, 10)

	return StoreClient{
		Conn: conn,
		URL:  urL,
		RPC:  clientRPC,
		PTZ:  clientPTZ,
		CGI:  clientCGI,
		File: clientFile,
	}
}

type StoreClient struct {
	Conn Conn
	URL  *url.URL
	RPC  dahuarpc.Client
	PTZ  ptz.Client
	CGI  dahuacgi.Client
	File dahuarpc.FileClient
}

func (c StoreClient) Close() error {
	return c.CloseContext(context.Background())
}

func (c StoreClient) CloseContext(ctx context.Context) error {
	c.File.Close()
	return c.RPC.Close(ctx)
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db:        db,
		clientsMu: sync.Mutex{},
		clients:   make(map[int64]StoreClient),
	}
}

// Store handles creating and caching clients.
type Store struct {
	db *sqlx.DB

	clientsMu sync.Mutex
	clients   map[int64]StoreClient
}

func (*Store) String() string {
	return "dahua.Store"
}

// Close closes all clients.
func (s *Store) Close() {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), StoreClientCloseTimeout)
	defer cancel()

	wg := sync.WaitGroup{}
	for _, client := range s.clients {
		wg.Add(1)
		go func(client StoreClient) {
			checkStoreClientCloseError(client.Conn, client.CloseContext(ctx))
			wg.Done()
		}(client)
	}
	wg.Wait()
}

func (s *Store) Serve(ctx context.Context) error {
	<-ctx.Done()
	s.Close()
	return ctx.Err()
}

func (s *Store) GetClient(ctx context.Context, deviceKey types.Key) (StoreClient, error) {
	s.clientsMu.Lock()
	var conn Conn
	err := s.db.GetContext(ctx, &conn, `
		SELECT id, uuid, name, ip, username, password
		FROM dahua_devices WHERE id = ? OR uuid = ?
	`, deviceKey.ID, deviceKey.UUID)
	if err != nil {
		s.clientsMu.Unlock()
		return StoreClient{}, err
	}

	client, ok := s.clients[conn.Key.ID]
	if !ok {
		// Not found

		client = NewStoreClient(conn)
		s.clients[conn.Key.ID] = client
	} else if !client.Conn.EQ(conn) {
		// Found but not equal

		go func(client StoreClient) { checkStoreClientCloseError(client.Conn, client.Close()) }(client)

		client = NewStoreClient(conn)
		s.clients[conn.Key.ID] = client
	}
	s.clientsMu.Unlock()

	return client, nil
}

func (s *Store) Register() *Store {
	bus.Subscribe(s.String(), func(ctx context.Context, event bus.DeviceDeleted) error {
		ctx, cancel := context.WithTimeout(ctx, StoreClientCloseTimeout)
		defer cancel()

		s.clientsMu.Lock()
		defer s.clientsMu.Unlock()

		wg := sync.WaitGroup{}
		for _, key := range event.DeviceKeys {
			wg.Add(1)

			client, found := s.clients[key.ID]
			if !found {
				continue
			}

			delete(s.clients, key.ID)

			go func(client StoreClient) {
				checkStoreClientCloseError(client.Conn, client.CloseContext(ctx))
				wg.Done()
			}(client)
		}
		wg.Wait()

		return nil
	})
	return s
}
