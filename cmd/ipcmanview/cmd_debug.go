package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
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
			conn := dahua.NewClient(device.DahuaConn)
			defer conn.RPC.Close(context.Background())
			defer wg.Done()

			cfg, err := config.GetVideoAnalyseRules(ctx, conn.RPC)
			if err != nil {
				log.Err(err).Send()
				return
			}

			_, err = json.MarshalIndent(cfg.Tables[0].Data, "", "  ")
			if err != nil {
				log.Err(err).Send()
				return
			}

			for i := range cfg.Tables[0].Data {
				fmt.Println(cfg.Tables[0].Data[i].Enable, cfg.Tables[0].Data[i].Name)
				cfg.Tables[0].Data[i].Enable = false
				// fmt.Println(string(b))
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
