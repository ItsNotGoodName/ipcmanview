package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
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

	for _, device := range devices {
		conn := dahua.NewClient(device.DahuaConn)
		defer conn.RPC.Close(context.Background())

		cfg, err := config.GetLocales(ctx, conn.RPC)
		if err != nil {
			log.Err(err).Str("name", device.Name).Send()
			continue
		}

		b, err := json.MarshalIndent(cfg.Tables[0].JSON, "", "  ")
		if err != nil {
			return err
		}

		fmt.Printf("name=%s, json=%s\n", device.Name, string(b))

		conn.RPC.Close(context.Background())
	}

	return nil
}
