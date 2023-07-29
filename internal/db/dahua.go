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

func DahuaCameraUpdate(ctx Context, r core.DahuaCameraUpdate) (core.DahuaCamera, error) {
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
	err := ScanOne(ctx.Context, ctx.Conn, &camera, DahuaCameras.
		UPDATE(cols).
		MODEL(model.DahuaCameras{
			Address:  r.DahuaCamera.Address,
			Username: r.DahuaCamera.Username,
			Password: r.DahuaCamera.Password,
		}).
		WHERE(DahuaCameras.ID.EQ(Int64(r.ID))).
		RETURNING(pDahuaCamera),
	)
	return camera, err
}
