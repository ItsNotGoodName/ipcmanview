package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	webdahua "github.com/ItsNotGoodName/ipcmanview/internal/web/dahua"
)

type CmdScan struct {
	Shared
	SharedCameras
	Full  bool `help:"Run full file scan."`
	Reset bool `help:"Reset all file cursors."`
}

func (c *CmdScan) Run(ctx *Context) error {
	db, err := useDB(ctx, c.DBPath)
	if err != nil {
		return err
	}

	cameras, err := c.useCameras(ctx, db)
	if err != nil {
		return err
	}

	scanType := webdahua.ScanTypeQuick
	if c.Full {
		scanType = webdahua.ScanTypeFull
	}

	if c.Reset {
		for _, camera := range cameras {
			webdahua.ScanReset(ctx, db, camera.ID)
		}
	}

	for _, camera := range cameras {
		conn := dahua.NewConn(camera.DahuaCamera)
		defer conn.RPC.Close(context.Background())

		err = webdahua.Scan(ctx, db, conn.RPC, conn.Camera, scanType)
		if err != nil {
			return err
		}

		conn.RPC.Close(context.Background())
	}

	return nil
}
