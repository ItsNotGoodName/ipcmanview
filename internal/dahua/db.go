package dahua

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/dbgen/postgres/dahua/model"
	dahua "github.com/ItsNotGoodName/ipcmango/internal/dbgen/postgres/dahua/table"
	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
	. "github.com/go-jet/jet/v2/postgres"
)

type dbT struct{}

var DB dbT

var dbCameraProjection ProjectionList = []Projection{
	dahua.Cameras.ID.AS("id"),
	dahua.Cameras.Address.AS("address"),
	dahua.Cameras.Username.AS("username"),
	dahua.Cameras.Password.AS("password"),
	dahua.Cameras.Location.AS("location"),
	dahua.Cameras.CreatedAt.AS("created_at"),
}

type dbCamera struct {
	model.Cameras
	Location models.Location
}

func (dbT) CameraCreate(ctx context.Context, db qes.Querier, r core.DahuaCamera) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		INSERT(dahua.Cameras.Address, dahua.Cameras.Username, dahua.Cameras.Password, dahua.Cameras.Location).
		MODEL(dbCamera{
			Cameras: model.Cameras{
				Address:  r.Address,
				Username: r.Username,
				Password: r.Password,
			},
			Location: r.Location,
		}).
		RETURNING(dbCameraProjection),
	)
	return camera, err
}

func (dbT) CameraUpdate(ctx context.Context, db qes.Querier, r *core.DahuaCameraUpdate) (core.DahuaCamera, error) {
	value, err := r.Value()
	if err != nil {
		return value, err
	}

	var cols ColumnList
	if r.Address {
		cols = append(cols, dahua.Cameras.Address)
	}
	if r.Username {
		cols = append(cols, dahua.Cameras.Username)
	}
	if r.Password {
		cols = append(cols, dahua.Cameras.Password)
	}

	var camera core.DahuaCamera
	err = qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		UPDATE(cols).
		MODEL(model.Cameras{
			Address:  value.Address,
			Username: value.Username,
			Password: value.Password,
		}).
		WHERE(dahua.Cameras.ID.EQ(Int64(value.ID))).
		RETURNING(dbCameraProjection),
	)
	return camera, err
}

func (dbT) CameraGet(ctx context.Context, db qes.Querier, id int64) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		SELECT(dbCameraProjection).
		WHERE(dahua.Cameras.ID.EQ(Int64(id))))
	return camera, err
}

func (dbT) CameraGetByAddress(ctx context.Context, db qes.Querier, address string) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		SELECT(dbCameraProjection).
		WHERE(dahua.Cameras.Address.EQ(String(address))))
	return camera, err
}

func (dbT) CameraDelete(ctx context.Context, db qes.Querier, id int64) error {
	_, err := qes.ExecOne(ctx, db, dahua.Cameras.
		DELETE().
		WHERE(dahua.Cameras.ID.EQ(Int64(id))).
		RETURNING(dahua.Cameras.ID))
	return err
}

func (dbT) ScanCameraFilesUpsert(ctx context.Context, db qes.Querier, cam models.DahuaScanCamera, files []mediafilefind.FindNextFileInfo, updatedAt time.Time) (int64, error) {
	if len(files) == 0 {
		return 0, nil
	}

	stmt := dahua.CameraFiles.INSERT(
		dahua.CameraFiles.CameraID,
		dahua.CameraFiles.FilePath,
		dahua.CameraFiles.Kind,
		dahua.CameraFiles.Size,
		dahua.CameraFiles.StartTime,
		dahua.CameraFiles.EndTime,
		dahua.CameraFiles.Events,
		dahua.CameraFiles.UpdatedAt,
	)

	for _, file := range files {
		startTime, endTime, err := file.UniqueTime(cam.Seed, cam.Location)
		if err != nil {
			return 0, err
		}

		events, err := json.Marshal(file.Events)
		if err != nil {
			return 0, err
		}

		stmt = stmt.MODEL(
			model.CameraFiles{
				CameraID:  int32(cam.ID),
				FilePath:  file.FilePath,
				Kind:      file.Type,
				Size:      int32(file.Length),
				StartTime: startTime,
				EndTime:   endTime,
				UpdatedAt: updatedAt,
				Events:    string(events),
			},
		)
	}

	stmt = stmt.
		ON_CONFLICT(
			dahua.CameraFiles.CameraID,
			dahua.CameraFiles.FilePath,
		).
		DO_UPDATE(SET(
			dahua.CameraFiles.UpdatedAt.SET(dahua.CameraFiles.EXCLUDED.UpdatedAt),
		)).
		ON_CONFLICT(dahua.CameraFiles.StartTime).
		DO_NOTHING()
	res, err := qes.Exec(ctx, db, stmt)
	return res.RowsAffected(), err
}

func (dbT) ScanCameraFilesDelete(ctx context.Context, db qes.Querier, cameraID int64, scanPeriod ScanPeriod, updatedAt time.Time) (int64, error) {
	stmt := dahua.CameraFiles.
		DELETE().
		WHERE(dahua.CameraFiles.UpdatedAt.LT(TimestampzT(updatedAt)).
			AND(dahua.CameraFiles.CameraID.EQ(Int64(cameraID)).
				AND(dahua.CameraFiles.StartTime.BETWEEN(TimestampzT(scanPeriod.Start), TimestampzT(scanPeriod.End)))))
	res, err := qes.Exec(ctx, db, stmt)
	return res.RowsAffected(), err
}
