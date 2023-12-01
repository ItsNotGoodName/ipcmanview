package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

type Context struct {
	context.Context
	Debug bool
}

var mainCmd struct {
	Debug bool `help:"Enable debug mode."`

	Web  WebCmd  `cmd:"" help:"Start web server."`
	Scan ScanCmd `cmd:"" help:"Scan files on cameras."`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ktx := kong.Parse(&mainCmd,
		kong.Description("Application for managing and viewing Dahua IP cameras."))
	err := ktx.Run(&Context{
		Context: ctx,
		Debug:   mainCmd.Debug,
	})
	ktx.FatalIfErrorf(err)
}
