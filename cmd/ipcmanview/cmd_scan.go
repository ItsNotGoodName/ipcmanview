package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
)

type CmdScan struct {
	Shared
	SharedDevices
	Full  bool `help:"Run full file scan."`
	Reset bool `help:"Reset all file cursors."`
}

func (c *CmdScan) Run(ctx *Context) error {
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	devices, err := c.useDevices(ctx, db)
	if err != nil {
		return err
	}

	scanType := dahua.ScanTypeQuick
	if c.Full {
		scanType = dahua.ScanTypeFull
	}

	for _, device := range devices {
		err := dahua.ScanLockCreate(ctx, db, device.DahuaConn.ID)
		if err != nil {
			return err
		}
		cancel := dahua.ScanLockHeartbeat(ctx, db, device.DahuaConn.ID)
		defer cancel()

		if c.Reset {
			err := dahua.ScanReset(ctx, db, device.DahuaDevice.ID)
			if err != nil {
				return err
			}
		}

		conn := dahuacore.NewConn(device.DahuaConn)
		defer conn.RPC.Close(context.Background())

		err = dahua.Scan(ctx, db, conn.RPC, conn.Device, scanType)
		if err != nil {
			return err
		}

		cancel()
		conn.RPC.Close(context.Background())
	}

	return nil
}
