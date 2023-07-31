package db

import (
	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/public/model"
	. "github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

var pDahuaCamera ProjectionList = []Projection{
	DahuaCameras.ID.AS("id"),
	DahuaCameras.Address.AS("address"),
	DahuaCameras.Username.AS("username"),
	DahuaCameras.Password.AS("password"),
	DahuaCameras.CreatedAt.AS("created_at"),
}

func DahuaCameraCreate(ctx Context, r core.DahuaCamera) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := ScanOne(ctx.Context, ctx.Conn, &camera, DahuaCameras.
		INSERT(DahuaCameras.Address, DahuaCameras.Username, DahuaCameras.Password).
		MODEL(model.DahuaCameras{
			Address:  r.Address,
			Username: r.Username,
			Password: r.Password,
		}).
		RETURNING(pDahuaCamera),
	)
	return camera, err
}

func DahuaCameraUpdate(ctx Context, r *core.DahuaCameraUpdate) (core.DahuaCamera, error) {
	value, err := r.Value()
	if err != nil {
		return value, err
	}

	var cols ColumnList
	if r.Address {
		cols = append(cols, DahuaCameras.Address)
	}
	if r.Username {
		cols = append(cols, DahuaCameras.Username)
	}
	if r.Password {
		cols = append(cols, DahuaCameras.Password)
	}

	var camera core.DahuaCamera
	err = ScanOne(ctx.Context, ctx.Conn, &camera, DahuaCameras.
		UPDATE(cols).
		MODEL(model.DahuaCameras{
			Address:  value.Address,
			Username: value.Username,
			Password: value.Password,
		}).
		WHERE(DahuaCameras.ID.EQ(Int64(value.ID))).
		RETURNING(pDahuaCamera),
	)
	return camera, err
}

func DahuaCameraGet(ctx Context, id int64) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := ScanOne(ctx.Context, ctx.Conn, &camera, DahuaCameras.
		SELECT(pDahuaCamera).
		WHERE(DahuaCameras.ID.EQ(Int64(id))))
	return camera, err
}

func DahuaCameraDelete(ctx Context, id int64) error {
	_, err := ExecOne(ctx.Context, ctx.Conn, DahuaCameras.DELETE().WHERE(DahuaCameras.ID.EQ(Int64(id))).RETURNING(DahuaCameras.ID))
	return err
}
