package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type Context struct {
	context.Context
	Debug bool
}

var mainCmd struct {
	Debug bool `env:"DEBUG" help:"Enable debug mode."`

	Version CmdVersion `cmd:"" help:"Show version."`
	Serve   CmdServe   `cmd:"" help:"Start application."`
	Debug_  CmdDebug   `name:"debug" cmd:""`
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
	// Get levels
	zerologLevel := zerolog.InfoLevel
	slogLevel := slog.LevelInfo
	if debug {
		zerologLevel = zerolog.DebugLevel
		slogLevel = slog.LevelDebug
	}

	// Set loggers
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerologLevel)
	slog.SetDefault(slog.New(slogzerolog.Option{Level: slogLevel, Logger: &log.Logger}.NewZerologHandler()))
}
