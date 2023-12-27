package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/rs/zerolog/log"
)

type CmdDebug struct {
	Shared
	SharedDevices
}

func (c *CmdDebug) Run(ctx *Context) error {
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	devices, err := c.useDevices(ctx, db)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, device := range devices {
		wg.Add(1)
		go func(device models.DahuaDeviceConn) {
			client := dahua.NewClient(device.DahuaConn)
			defer client.RPC.Close(context.Background())
			defer wg.Done()

			status, err := ptz.GetStatus(ctx, client.PTZ, 0)
			if err != nil {
				log.Err(err).Send()
				return
			}

			fmt.Println(status)
		}(device)
	}
	wg.Wait()

	return nil
}
