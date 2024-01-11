package main

import (
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
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

	dahuaFilesFS, err := c.useDahuaFileFS()
	if err != nil {
		return err
	}

	deleted, err := dahua.DeleteAllOrphanAferoFile(ctx, db, dahuaFilesFS)
	if err != nil {
		return err
	}

	fmt.Println(deleted)

	// for _, device := range devices {
	// 	conn := dahua.NewClient(device.DahuaConn)
	// 	defer conn.RPC.Close(context.Background())
	//
	// 	cfg, err := config.GetRecord(ctx, conn.RPC)
	// 	if err != nil {
	// 		log.Err(err).Str("name", device.Name).Send()
	// 		continue
	// 	}
	//
	// 	b, err := json.MarshalIndent(cfg.Tables[0].Data.TimeSection, "", "  ")
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	fmt.Printf("name=%s, json=%s\n", device.Name, string(b))
	//
	// 	conn.RPC.Close(context.Background())
	// }

	return nil
}
