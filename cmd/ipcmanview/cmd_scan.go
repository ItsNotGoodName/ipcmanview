package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
)

type ScanCmd struct {
	DBPath string `env:"DB_PATH" default:"sqlite.db"`
	Full   bool
}

func (c *ScanCmd) Run(ctx *Context) error {
	db, err := useDB(c.DBPath)
	if err != nil {
		return err
	}

	dbCamera, err := db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	scanType := webdahua.ScanTypeQuick
	if c.Full {
		scanType = webdahua.ScanTypeFull
	}

	for _, camera := range webdahua.ConvertListDahuaCameraRows(dbCamera) {
		conn := dahua.NewConn(camera)
		defer conn.RPC.Close(context.Background())

		err = webdahua.Scan(ctx, db, conn.RPC, conn.Camera, scanType)
		if err != nil {
			return err
		}

		conn.RPC.Close(context.Background())
	}

	return nil
}
