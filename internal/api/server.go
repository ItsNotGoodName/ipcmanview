package api

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
)

type DahuaConnStore interface {
	Get(ctx context.Context, id int64) (models.DahuaConn, bool, error)
	List(ctx context.Context) ([]models.DahuaConn, error)
}

func NewServer(pub pubsub.Pub, dahuaStore *dahua.Store, dahuaConnStore DahuaConnStore) *Server {
	return &Server{
		pub:            pub,
		dahuaStore:     dahuaStore,
		dahuaConnStore: dahuaConnStore,
	}
}

type Server struct {
	pub            pubsub.Pub
	dahuaStore     *dahua.Store
	dahuaConnStore DahuaConnStore
}
