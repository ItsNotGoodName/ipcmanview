package main

import (
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/gorise"
)

type CmdDebug struct {
	Shared
}

func (c *CmdDebug) Run(ctx *Context) error {
	if err := c.init(); err != nil {
		return err
	}

	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	hub := bus.NewHub(ctx)

	afs, err := c.useDahuaAFS()
	if err != nil {
		return err
	}

	store := dahua.NewStore().Register(hub)

	scanLocker := dahua.NewScanLocker()

	dahua.Init(dahua.App{
		DB:         db,
		Hub:        hub,
		AFS:        afs,
		Store:      store,
		ScanLocker: scanLocker,
	})

	urL, _ := os.LookupEnv("SENDER_URL")

	sender, err := gorise.Build(urL)
	if err != nil {
		return err
	}

	f, err := afs.Open("8154a5ae-fcfd-41e4-be66-edacea89d255.jpg")
	if err != nil {
		return err
	}
	defer f.Close()

	return sender.Send(ctx, gorise.Message{
		Title: "Test title",
		Body:  "Test body.",
		Attachments: []gorise.Attachment{
			{
				Name:   "Test",
				Mime:   "image/jpeg",
				Reader: f,
			},
		},
	})
}

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
