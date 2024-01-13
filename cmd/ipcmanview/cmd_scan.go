package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	dahua1 "github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
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

	scanType := models.DahuaScanTypeQuick
	if c.Full {
		scanType = models.DahuaScanTypeFull
	}

	for _, device := range devices {
		if c.Reset {
			err := dahua.ScanReset(ctx, db, device.DahuaDevice.ID)
			if err != nil {
				return err
			}
		}

		conn := dahua1.NewClient(device.DahuaConn)
		defer conn.RPC.Close(context.Background())

		err = dahua.Scan(ctx, db, conn.RPC, conn.Conn, scanType)
		if err != nil {
			return err
		}

		conn.RPC.Close(context.Background())
	}

	return nil
}
