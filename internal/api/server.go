package api

import (
	"context"
	"io"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
)

type DahuaRepo interface {
	GetConn(ctx context.Context, id int64) (models.DahuaConn, bool, error)
	ListConn(ctx context.Context) ([]models.DahuaConn, error)
	GetFileByFilePath(ctx context.Context, cameraID int64, filePath string) (models.DahuaFile, error)
}

type DahuaFileCache interface {
	// Save(ctx context.Context, file models.DahuaFile, r io.ReadCloser) error
	Exists(ctx context.Context, file models.DahuaFile) (bool, error)
	Get(ctx context.Context, file models.DahuaFile) (io.ReadCloser, error)
}

func NewServer(
	pub pubsub.Pub,
	dahuaStore *dahuacore.Store,
	dahuaRepo DahuaRepo,
	dahuaFileCache DahuaFileCache,
) *Server {
	return &Server{
		pub:            pub,
		dahuaStore:     dahuaStore,
		dahuaRepo:      dahuaRepo,
		dahuaFileCache: dahuaFileCache,
	}
}

type Server struct {
	pub            pubsub.Pub
	dahuaStore     *dahuacore.Store
	dahuaRepo      DahuaRepo
	dahuaFileCache DahuaFileCache
}
