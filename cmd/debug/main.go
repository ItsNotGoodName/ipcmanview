package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/db"
	"github.com/ItsNotGoodName/ipcmanview/migrations"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/magicbox"
	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx, shutdown := interrupt.Context()
	defer shutdown()

	// Database
	pool, err := db.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database connection pool")
	}
	defer pool.Close()

	// Database migrate
	if err := migrations.Migrate(ctx, pool); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	bus := db.NewBusFromPool(pool)
	super.Add(bus)

	dahuaSuper := dahua.NewSupervisor(pool)
	dahuaSuper.Register(bus)
	super.Add(dahuaSuper)

	superDoneC := super.ServeBackground(ctx)

	// --------------------------------------------------------------------------
	{
		seed(ctx, pool)
	}

	// --------------------------------------------------------------------------
	super.Add(sutureext.NewServiceFunc("debug", func(ctx context.Context) error {
		time.Sleep(1 * time.Second)
		c, err := dahuaSuper.GetOrCreateWorker(ctx, 319)
		if err != nil {
			return err
		}

		sn, err := magicbox.GetSerialNo(ctx, c)
		if err != nil {
			return err
		}

		fmt.Println("SN:", sn)

		return suture.ErrDoNotRestart
	}))

	if err := <-superDoneC; err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
