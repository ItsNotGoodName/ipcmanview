package dahua

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/auth"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
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

type StoreHooks interface {
	ConnCreated(conn Conn)
	ConnUpdated(conn Conn)
}

type Store struct {
	clientsMu sync.Mutex
	clients   map[string]storeClient
	hooks     StoreHooks
}

func NewStore(hooks StoreHooks) *Store {
	return &Store{
		clientsMu: sync.Mutex{},
		clients:   make(map[string]storeClient),
		hooks:     hooks,
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
	created := false
	updated := false
	if !ok {
		client = newStoreClient(camera)
		s.clients[camera.ID] = client
		created = true
	} else if !cameraIdentityEqual(client.Camera, camera) {
		client.ConnRPC.Close(ctx)
		client = newStoreClient(camera)
		s.clients[camera.ID] = client
		updated = true
	}

	conn := newConn(client)

	if created {
		s.hooks.ConnCreated(conn)
	}
	if updated {
		s.hooks.ConnUpdated(conn)
	}

	return conn
}

func (s *Store) ConnByCamera(ctx context.Context, camera models.DahuaCamera) Conn {
	s.clientsMu.Lock()
	conn := s.getOrCreateCamera(ctx, camera)
	s.clientsMu.Unlock()
	return conn
}

func (s *Store) ConnListByCameras(ctx context.Context, cameras ...models.DahuaCamera) ([]Conn, error) {
	s.clientsMu.Lock()
	for _, camera := range cameras {
		_ = s.getOrCreateCamera(ctx, camera)
	}

	clients := make([]Conn, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, newConn(client))
	}
	s.clientsMu.Unlock()

	slices.SortFunc(clients, func(a, b Conn) int { return cmp.Compare(a.Camera.ID, b.Camera.ID) })

	return clients, nil
}

func (s *Store) ConnByID(ctx context.Context, id string) (Conn, error) {
	s.clientsMu.Lock()
	client, ok := s.clients[id]
	s.clientsMu.Unlock()
	if !ok {
		return Conn{}, fmt.Errorf("Conn not found by ID: %s", id)
	}

	return newConn(client), nil
}
