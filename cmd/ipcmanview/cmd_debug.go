package main

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
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
			conn := dahuacore.NewConn(device.DahuaConn)
			defer conn.RPC.Close(context.Background())
			defer wg.Done()

			config, err := configmanager.GetConfig[[]configmanager.VideoInMode](ctx, conn.RPC, "VideoInMode")
			if err != nil {
				log.Err(err).Send()
				return
			}

			err = configmanager.SetConfig(ctx, conn.RPC, "VideoInMode", config)
			if err != nil {
				log.Err(err).Send()
				return
			}
		}(device)
	}
	wg.Wait()

	return nil
}
