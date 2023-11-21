package webcore

import (
	"context"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlc"
)

type DahuaBus interface {
	OnCameraEvent(h func(ctx context.Context, evt models.EventDahuaCameraEvent) error)
}

func RegisterDahuaBus(bus DahuaBus, sqlcDB *sqlc.Queries) {
	bus.OnCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		cameraID, err := strconv.ParseInt(evt.Event.CameraID, 10, 64)
		if err != nil {
			return err
		}

		_, err = sqlcDB.CreateDahuaEvent(ctx, sqlc.CreateDahuaEventParams{
			CameraID:      cameraID,
			ContentType:   evt.Event.ContentType,
			ContentLength: int64(evt.Event.ContentLength),
			Code:          evt.Event.Code,
			Action:        evt.Event.Action,
			Index:         int64(evt.Event.Index),
			Data:          evt.Event.Data,
			CreatedAt:     evt.Event.CreatedAt,
		})
		return err
	})
}

func SyncDahuaStore(ctx context.Context, sqlcDB *sqlc.Queries, store *dahua.Store) error {
	dbCameras, err := sqlcDB.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}
	_, err = store.ConnListByCameras(ctx, ConvertDahuaCameras(dbCameras)...)
	return err
}

func ConvertDahuaCameras(dbCameras []sqlc.DahuaCamera) []models.DahuaCamera {
	cameras := make([]models.DahuaCamera, 0, len(dbCameras))
	for _, dbCamera := range dbCameras {
		cameras = append(cameras, ConvertDahuaCamera(dbCamera))
	}

	return cameras
}

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

func NewDahuaStoreProxy(dahuaStore *dahua.Store, sqlcDB *sqlc.Queries) DahuaStoreProxy {
	return DahuaStoreProxy{
		store: dahuaStore,
		db:    sqlcDB,
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
