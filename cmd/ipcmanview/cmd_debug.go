package main

import (
	"fmt"
	"net"
)

type CmdDebug struct {
	Shared
	SharedDevices
}

func (c *CmdDebug) Run(ctx *Context) error {
	// db, err := c.useDB(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// devices, err := c.useDevices(ctx, db)
	// if err != nil {
	// 	return err
	// }

	ip := "lookup"
	if net.IP(ip).To4() == nil {
		ips, err := net.LookupIP(ip)
		if err == nil {
			return err
		}

		for _, i2 := range ips {
			if i2.To4() != nil {
				ip = i2.String()
				break
			}
		}
	}
	fmt.Println(ip)

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
