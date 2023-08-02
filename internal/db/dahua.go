package db

import (
	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/dahua/model"
	dahua "github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/dahua/table"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
	. "github.com/go-jet/jet/v2/postgres"
)

var pDahuaCamera ProjectionList = []Projection{
	dahua.Cameras.ID.AS("id"),
	dahua.Cameras.Address.AS("address"),
	dahua.Cameras.Username.AS("username"),
	dahua.Cameras.Password.AS("password"),
	dahua.Cameras.CreatedAt.AS("created_at"),
}

func DahuaCameraCreate(dbCtx Context, r core.DahuaCamera) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := qes.ScanOne(dbCtx.Context, dbCtx.Conn, &camera, dahua.Cameras.
		INSERT(dahua.Cameras.Address, dahua.Cameras.Username, dahua.Cameras.Password).
		MODEL(model.Cameras{
			Address:  r.Address,
			Username: r.Username,
			Password: r.Password,
		}).
		RETURNING(pDahuaCamera),
	)
	return camera, err
}

func DahuaCameraUpdate(dbCtx Context, r *core.DahuaCameraUpdate) (core.DahuaCamera, error) {
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
	err = qes.ScanOne(dbCtx.Context, dbCtx.Conn, &camera, dahua.Cameras.
		UPDATE(cols).
		MODEL(model.Cameras{
			Address:  value.Address,
			Username: value.Username,
			Password: value.Password,
		}).
		WHERE(dahua.Cameras.ID.EQ(Int64(value.ID))).
		RETURNING(pDahuaCamera),
	)
	return camera, err
}

func DahuaCameraGet(dbCtx Context, id int64) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := qes.ScanOne(dbCtx.Context, dbCtx.Conn, &camera, dahua.Cameras.
		SELECT(pDahuaCamera).
		WHERE(dahua.Cameras.ID.EQ(Int64(id))))
	return camera, err
}

func DahuaCameraGetByAddress(ctx Context, address string) (core.DahuaCamera, error) {
	var camera core.DahuaCamera
	err := qes.ScanOne(ctx.Context, ctx.Conn, &camera, dahua.Cameras.
		SELECT(pDahuaCamera).
		WHERE(dahua.Cameras.Address.EQ(String(address))))
	return camera, err
}

func DahuaCameraDelete(dbCtx Context, id int64) error {
	_, err := qes.ExecOne(dbCtx.Context, dbCtx.Conn, dahua.Cameras.DELETE().WHERE(dahua.Cameras.ID.EQ(Int64(id))).RETURNING(dahua.Cameras.ID))
	return err
}
