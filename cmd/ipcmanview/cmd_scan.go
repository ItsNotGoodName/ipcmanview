package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuaweb"
)

type CmdScan struct {
	Shared
	SharedCameras
	Full  bool `help:"Run full file scan."`
	Reset bool `help:"Reset all file cursors."`
}

func (c *CmdScan) Run(ctx *Context) error {
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	cameras, err := c.useCameras(ctx, db)
	if err != nil {
		return err
	}

	scanType := dahuaweb.ScanTypeQuick
	if c.Full {
		scanType = dahuaweb.ScanTypeFull
	}

	if c.Reset {
		for _, camera := range cameras {
			dahuaweb.ScanReset(ctx, db, camera.DahuaCamera.ID)
		}
	}

	for _, camera := range cameras {
		conn := dahua.NewConn(camera.DahuaConn)
		defer conn.RPC.Close(context.Background())

		err = dahuaweb.Scan(ctx, db, conn.RPC, conn.Camera, scanType)
		if err != nil {
			return err
		}

		conn.RPC.Close(context.Background())
	}

	return nil
}
