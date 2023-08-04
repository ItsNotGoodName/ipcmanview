package main

import (
	"context"
	"os"

	"github.com/ItsNotGoodName/ipcmango/internal/build"
	"github.com/ItsNotGoodName/ipcmango/internal/dahua"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/internal/event"
	"github.com/ItsNotGoodName/ipcmango/migrations"
	"github.com/ItsNotGoodName/ipcmango/pkg/background"
	"github.com/ItsNotGoodName/ipcmango/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmango/sandbox"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// sandbox.Dahua(interrupt.Context())
	// return

	ctx, shutdown := context.WithCancel(interrupt.Context())
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

	sandbox.Sandbox(ctx, pool)
	return

	// Event Bus
	bus := event.NewBus()

	// Dahua
	store := dahua.NewStore()
	dahua.StoreConnectBus(bus, store, pool)

	// sandbox.Jet(ctx, pool)
	<-background.Run(ctx,
		sandbox.Chi(ctx, shutdown, pool),
		event.Background(bus, pool),
		// background.NewFunction(background.BlockingForever, func(ctx context.Context) {
		// 	pool.AcquireFunc(ctx, func(c *pgxpool.Conn) error {
		// 		username, _ := os.LookupEnv("IPC_USERNAME")
		// 		password, _ := os.LookupEnv("IPC_PASSWORD")
		// 		ip, _ := os.LookupEnv("IPC_IP")
		//
		// 		dbCtx := db.Conn{Context: ctx, Conn: c.Conn()}
		// 		var cam core.DahuaCamera
		// 		{
		// 			req, err := core.NewDahuaCamera(core.DahuaCameraCreate{
		// 				Address:  ip,
		// 				Username: username,
		// 				Password: password,
		// 			})
		// 			if err != nil {
		// 				if errors.Is(err, context.Canceled) {
		// 					return nil
		// 				}
		// 				panic(err)
		// 			}
		//
		// 			ctx := dahua.DB(dbCtx)
		//
		// 			cam, err = ctx.CameraCreate(req)
		// 			if err != nil {
		// 				cam, err = ctx.CameraGetByAddress(req.Address)
		// 				if err != nil {
		// 					return err
		// 				}
		// 			}
		// 		}
		//
		// 		for {
		// 			log.Debug().Msg("Sleeping...")
		// 			time.Sleep(1 * time.Second)
		//
		// 			actor, err := store.GetOrCreate(dahua.DB(dbCtx), cam.ID)
		// 			if err != nil {
		// 				if errors.Is(err, context.Canceled) {
		// 					return nil
		// 				}
		// 				log.Err(err).Msg("Failed to get actor")
		// 				continue
		// 			}
		//
		// 			sn, err := magicbox.GetSerialNo(ctx, actor)
		// 			if err != nil {
		// 				if errors.Is(err, context.Canceled) {
		// 					return nil
		// 				}
		// 				log.Err(err).Msg("Failed to get sn")
		// 				continue
		// 			}
		// 			fmt.Println(sn)
		// 		}
		// 	})
		// }),
	)
	// sandbox.User(ctx, pool)
}

var (
	builtBy    = "unknown"
	commit     = ""
	date       = ""
	version    = "dev"
	repoURL    = "https://github.com/ItsNotGoodName/smtpbridge"
	releaseURL = ""
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	build.Current = build.Build{
		BuiltBy:    builtBy,
		Commit:     commit,
		Date:       date,
		Version:    version,
		RepoURL:    repoURL,
		ReleaseURL: releaseURL,
	}
}
