package webcore

import (
	"context"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
)

func ConvertDahuaCamera(c sqlc.DahuaCamera) models.DahuaCamera {
	return models.DahuaCamera{
		ID:        strconv.FormatInt(c.ID, 10),
		Address:   c.Address,
		Username:  c.Username,
		Password:  c.Password,
		Location:  c.Location,
		Seed:      0,
		CreatedAt: c.CreatedAt,
	}
}

func NewDahuaStoreProxy(dahuaStore *dahua.Store, db *sqlc.Queries) DahuaStoreProxy {
	return DahuaStoreProxy{
		store: dahuaStore,
		db:    db,
	}
}

type DahuaStoreProxy struct {
	store *dahua.Store
	db    *sqlc.Queries
}

func (s DahuaStoreProxy) ConnByID(ctx context.Context, idStr string) (dahua.Conn, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return dahua.Conn{}, err
	}

	dbCamera, err := s.db.GetDahuaCamera(ctx, id)
	if err != nil {
		return dahua.Conn{}, err
	}

	return s.store.ConnByCamera(ctx, ConvertDahuaCamera(dbCamera)), nil
}

func (s DahuaStoreProxy) ConnListByCameras(ctx context.Context, _ ...models.DahuaCamera) ([]dahua.Conn, error) {
	dbCameras, err := s.db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	cameras := make([]models.DahuaCamera, 0, len(dbCameras))
	for _, dbCamera := range dbCameras {
		cameras = append(cameras, ConvertDahuaCamera(dbCamera))
	}

	conns, err := s.store.ConnListByCameras(ctx, cameras...)
	if err != nil {
		return nil, err
	}

	return conns, nil
}
