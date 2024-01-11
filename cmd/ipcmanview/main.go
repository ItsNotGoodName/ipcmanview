package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Context struct {
	context.Context
	Debug bool
}

var mainCmd struct {
	Debug bool `env:"DEBUG" help:"Enable debug mode."`

	Version CmdVersion `cmd:"" help:"Show version."`
	Serve   CmdServe   `cmd:"" help:"Start application."`
	Scan    CmdScan    `cmd:"" help:"Scan files on devices."`
	RPC     CmdRPC     `cmd:"" help:"Run RPC on devices."`
	Debug_  CmdDebug   `name:"debug" cmd:"" help:"Debug."`
}

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal().Err(err).Send()
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ktx := kong.Parse(&mainCmd, kong.Description("Application for managing and viewing Dahua devices."))

	initLogger(mainCmd.Debug)

	err = ktx.Run(&Context{
		Context: ctx,
		Debug:   mainCmd.Debug,
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal().Err(err).Send()
	}
}

func initLogger(debug bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if debug {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}
}
