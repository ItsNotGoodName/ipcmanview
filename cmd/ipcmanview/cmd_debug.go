package main

import (
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/endpoint"
)

type CmdDebug struct {
	Shared
}

func (c *CmdDebug) Run(ctx *Context) error {
	urL, _ := os.LookupEnv("SENDER_URL")

	sender, err := endpoint.Build(urL)
	if err != nil {
		return err
	}

	return sender.Send(ctx, endpoint.Message{
		Title:       "Test title",
		Body:        "Test body.",
		Attachments: []endpoint.Attachment{},
	})
}

// func (c *CmdDebug) Run(ctx *Context) error {
// 	if err := c.init(); err != nil {
// 		return err
// 	}
//
// 	db, err := c.useDB(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	hub := bus.NewHub(ctx)
//
// 	dahua.Init(dahua.App{
// 		DB:         db,
// 		Hub:        hub,
// 		AFS:        nil,
// 		Store:      nil,
// 		ScanLocker: dahua.ScanLocker{},
// 	})
//
// 	hub.OnDahuaFileCreated("DEBUG", func(ctx context.Context, event bus.DahuaFileCreated) error {
// 		fmt.Println("DEVICE:", event.DeviceID, "COUNT", event.Count)
// 		return nil
// 	})
//
// 	store := dahua.NewStore()
//
// 	start := time.Now()
//
// 	conns, err := store.ListClient(ctx)
//
// 	var wg sync.WaitGroup
//
// 	for _, c := range conns {
// 		wg.Add(1)
// 		go func(c dahua.Client) {
// 			defer wg.Done()
// 			err := dahua.Scan(ctx, c.RPC, c.Conn, models.DahuaScanType_Full)
// 			if err != nil {
// 				log.Err(err).Send()
// 			}
// 		}(c)
// 	}
//
// 	wg.Wait()
//
// 	fmt.Println("DURATION:", time.Now().Sub(start))
//
// 	return nil
// }
