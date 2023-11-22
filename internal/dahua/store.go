package dahua

import (
	"context"
	"fmt"
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

func cameraIdentityEqual(lhs, rhs models.DahuaCamera) bool {
	return lhs.Address == rhs.Address && lhs.Username == rhs.Username && lhs.Password == rhs.Password
}

func newStoreClient(camera models.DahuaCamera) storeClient {
	address := NewHTTPAddress(camera.Address)
	rpcHTTPClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	cgiHTTPClient := http.Client{}

	connRPC := auth.NewConn(dahuarpc.NewConn(rpcHTTPClient, address), camera.Username, camera.Password)
	connPTZ := ptz.NewClient(connRPC)
	connCGI := dahuacgi.NewConn(cgiHTTPClient, address, camera.Username, camera.Password)

	return storeClient{
		Camera:  camera,
		ConnRPC: connRPC,
		ConnPTZ: connPTZ,
		ConnCGI: connCGI,
	}
}

type storeClient struct {
	Camera  models.DahuaCamera
	ConnRPC *auth.Conn
	ConnPTZ *ptz.Client
	ConnCGI dahuacgi.Conn
}

func (c storeClient) Close(ctx context.Context) {
	if err := c.ConnRPC.Close(ctx); err != nil {
		log.Err(err).Int64("id", c.Camera.ID).Caller().Msg("Failed to close RPC connection")
	}
}

type StoreCameraStore interface {
	List(ctx context.Context) ([]models.DahuaCamera, error)
	Get(ctx context.Context, id int64) (models.DahuaCamera, bool, error)
}

type Store struct {
	cameraStore StoreCameraStore

	clientsMu sync.Mutex
	clients   map[int64]storeClient
}

func NewStore(cameraStore StoreCameraStore) *Store {
	return &Store{
		cameraStore: cameraStore,
		clientsMu:   sync.Mutex{},
		clients:     make(map[int64]storeClient),
	}
}

func (s *Store) Serve(ctx context.Context) error {
	<-ctx.Done()

	s.clientsMu.Lock()
	for _, rpcClient := range s.clients {
		rpcClient.ConnRPC.Close(context.TODO())
	}
	s.clientsMu.Unlock()

	return ctx.Err()
}

func (s *Store) getOrCreateCamera(ctx context.Context, camera models.DahuaCamera) Conn {
	client, ok := s.clients[camera.ID]
	if !ok {
		// Not found
		client = newStoreClient(camera)
		s.clients[camera.ID] = client
	} else if !cameraIdentityEqual(client.Camera, camera) {
		// Found but not equal

		client.Close(ctx)
		client = newStoreClient(camera)
		s.clients[camera.ID] = client
	}

	conn := newConn(client)

	return conn
}

func (s *Store) ConnList(ctx context.Context) ([]Conn, error) {
	s.clientsMu.Lock()
	cameras, err := s.cameraStore.List(ctx)
	if err != nil {
		s.clientsMu.Unlock()
		return nil, err
	}

	clients := make([]Conn, 0, len(s.clients))
	for _, camera := range cameras {
		clients = append(clients, s.getOrCreateCamera(ctx, camera))
	}
	s.clientsMu.Unlock()

	return clients, nil
}

func (s *Store) Conn(ctx context.Context, id int64) (Conn, error) {
	s.clientsMu.Lock()
	camera, found, err := s.cameraStore.Get(ctx, id)
	if err != nil {
		s.clientsMu.Unlock()
		return Conn{}, err
	}
	if !found {
		client, found := s.clients[id]
		if found {
			client.Close(ctx)
		}
		s.clientsMu.Unlock()
		return Conn{}, fmt.Errorf("Conn not found by ID: %d", id)
	}
	client := s.getOrCreateCamera(ctx, camera)
	s.clientsMu.Unlock()

	return client, nil
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
