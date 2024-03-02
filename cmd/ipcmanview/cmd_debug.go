package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/rs/zerolog/log"
)

type CmdDebug struct {
	Shared
}

func (c *CmdDebug) Run(ctx *Context) error {
	db, err := c.useDB(ctx)
	if err != nil {
		return err
	}

	bus := event.NewBus(ctx)

	bus.OnDahuaFileCreated("DEBUG", func(ctx context.Context, evt event.DahuaFileCreated) error {
		fmt.Println("DEVICE:", evt.DeviceID, "COUNT", evt.Count)
		return nil
	})

	store := dahua.NewStore(db)

	start := time.Now()

	conns, err := store.ListClient(ctx)

	var wg sync.WaitGroup

	for _, c := range conns {
		wg.Add(1)
		go func(c dahua.Client) {
			defer wg.Done()
			err := dahua.Scan(ctx, db, bus, c.RPC, c.Conn, models.DahuaScanTypeFull)
			if err != nil {
				log.Err(err).Send()
			}
		}(c)
	}

	wg.Wait()

	fmt.Println("DURATION:", time.Now().Sub(start))

	return nil
}
