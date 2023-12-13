package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
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

	scanType := dahua.ScanTypeQuick
	if c.Full {
		scanType = dahua.ScanTypeFull
	}

	if c.Reset {
		for _, camera := range cameras {
			err := dahua.ScanReset(ctx, db, camera.DahuaCamera.ID)
			if err != nil {
				return err
			}
		}
	}

	for _, camera := range cameras {
		conn := dahuacore.NewConn(camera.DahuaConn)
		defer conn.RPC.Close(context.Background())

		err = dahua.Scan(ctx, db, conn.RPC, conn.Camera, scanType)
		if err != nil {
			return err
		}

		conn.RPC.Close(context.Background())
	}

	return nil
}
