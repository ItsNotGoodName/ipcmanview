package dahua

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/auth"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/rs/zerolog/log"
)

func cameraEqual(lhs, rhs models.DahuaCamera) bool {
	return lhs.Address == rhs.Address && lhs.Username == rhs.Username && lhs.Password == rhs.Password
}

func newStoreClient(camera models.DahuaCamera, lastAccessed time.Time) storeClient {
	address := NewHTTPAddress(camera.Address)
	rpcHTTPClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	cgiHTTPClient := http.Client{}

	connRPC := auth.NewConn(dahuarpc.NewConn(rpcHTTPClient, address), camera.Username, camera.Password)
	connPTZ := ptz.NewClient(connRPC)
	connCGI := dahuacgi.NewConn(cgiHTTPClient, address, camera.Username, camera.Password)

	return storeClient{
		LastAccessed: lastAccessed,
		Camera:       camera,
		ConnRPC:      connRPC,
		ConnPTZ:      connPTZ,
		ConnCGI:      connCGI,
	}
}

type storeClient struct {
	LastAccessed time.Time
	Camera       models.DahuaCamera
	ConnRPC      *auth.Conn
	ConnPTZ      *ptz.Client
	ConnCGI      dahuacgi.Conn
}

func (c storeClient) Close(ctx context.Context) {
	if err := c.ConnRPC.Close(ctx); err != nil {
		log.Err(err).Int64("id", c.Camera.ID).Caller().Msg("Failed to close RPC connection")
	}
}

func newConn(c storeClient) Conn {
	return Conn{
		Camera: c.Camera,
		RPC:    c.ConnRPC,
		PTZ:    c.ConnPTZ,
		CGI:    c.ConnCGI,
	}
}

type Conn struct {
	Camera models.DahuaCamera
	RPC    *auth.Conn
	PTZ    *ptz.Client
	CGI    dahuacgi.Conn
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
			s.clientsMu.Lock()
			for _, rpcClient := range s.clients {
				rpcClient.Close(context.Background())
			}
			s.clientsMu.Unlock()

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

func (s *Store) getOrCreateCamera(ctx context.Context, camera models.DahuaCamera) Conn {
	client, ok := s.clients[camera.ID]
	if !ok {
		// Not found

		client = newStoreClient(camera, time.Now())
		s.clients[camera.ID] = client
	} else if !cameraEqual(client.Camera, camera) {
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

	conn := newConn(client)

	return conn
}

func (s *Store) ConnList(ctx context.Context, cameras []models.DahuaCamera) []Conn {
	clients := make([]Conn, 0, len(cameras))

	s.clientsMu.Lock()
	for _, camera := range cameras {
		clients = append(clients, s.getOrCreateCamera(ctx, camera))
	}
	s.clientsMu.Unlock()

	return clients
}

func (s *Store) Conn(ctx context.Context, camera models.DahuaCamera) Conn {
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

type StoreBus interface {
	OnCameraDeleted(h func(ctx context.Context, evt models.EventDahuaCameraDeleted) error)
}

func RegisterStoreBus(bus StoreBus, store *Store) {
	bus.OnCameraDeleted(func(ctx context.Context, evt models.EventDahuaCameraDeleted) error {
		store.ConnDelete(ctx, evt.CameraID)
		return nil
	})
}
