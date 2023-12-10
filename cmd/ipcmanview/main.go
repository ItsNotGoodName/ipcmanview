package main

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Context struct {
	context.Context
	Debug bool
}

var mainCmd struct {
	Debug bool `help:"Enable debug mode."`

	Serve CmdServe `cmd:"" help:"Start application."`
	Scan  CmdScan  `cmd:"" help:"Scan files on cameras."`
	RPC   CmdRPC   `cmd:"" help:"Run RPC on cameras."`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ktx := kong.Parse(&mainCmd, kong.Description("Application for managing and viewing Dahua IP cameras."))

	initLogger(mainCmd.Debug)

	err := ktx.Run(&Context{
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
