package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
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

			cfg, err := config.GetVideoAnalyseRules(ctx, conn.RPC)
			if err != nil {
				log.Err(err).Send()
				return
			}

			b, err := json.MarshalIndent(cfg.Tables[0].Data, "", "  ")
			if err != nil {
				log.Err(err).Send()
				return
			}
			fmt.Println(string(b))

			for i := range cfg.Tables[0].Data {
				if cfg.Tables[0].Data[i].Name == "IVS-1" {
					cfg.Tables[0].Data[i].Enable = false
				}
			}

			err = configmanager.SetConfig(ctx, conn.RPC, cfg)
			if err != nil {
				log.Err(err).Send()
				return
			}

		}(device)
	}
	wg.Wait()

	return nil
}
