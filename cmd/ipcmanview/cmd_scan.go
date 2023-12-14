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

	for _, camera := range cameras {
		err := dahua.ScanLockCreate(ctx, db, camera.DahuaConn.ID)
		if err != nil {
			return err
		}
		cancel := dahua.ScanLockHeartbeat(ctx, db, camera.DahuaConn.ID)
		defer cancel()

		if c.Reset {
			err := dahua.ScanReset(ctx, db, camera.DahuaCamera.ID)
			if err != nil {
				return err
			}
		}

		conn := dahuacore.NewConn(camera.DahuaConn)
		defer conn.RPC.Close(context.Background())

		err = dahua.Scan(ctx, db, conn.RPC, conn.Camera, scanType)
		if err != nil {
			return err
		}

		cancel()
		conn.RPC.Close(context.Background())
	}

	return nil
}
