package api

import (
	"context"
	"io"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
)

type DahuaConnStore interface {
	Get(ctx context.Context, id int64) (models.DahuaConn, bool, error)
	List(ctx context.Context) ([]models.DahuaConn, error)
}

type DahuaFileCache interface {
	// Save(ctx context.Context, file repo.DahuaFile, r io.ReadCloser) error
	Exists(ctx context.Context, file repo.DahuaFile) (bool, error)
	Get(ctx context.Context, file repo.DahuaFile) (io.ReadCloser, error)
}

func NewServer(pub pubsub.Pub, db repo.DB, dahuaStore *dahua.Store, dahuaConnStore DahuaConnStore, dahuaFileCache DahuaFileCache) *Server {
	return &Server{
		pub:            pub,
		db:             db,
		dahuaStore:     dahuaStore,
		dahuaConnStore: dahuaConnStore,
		dahuaFileCache: dahuaFileCache,
	}
}

type Server struct {
	pub            pubsub.Pub
	db             repo.DB
	dahuaStore     *dahua.Store
	dahuaConnStore DahuaConnStore
	dahuaFileCache DahuaFileCache
}
