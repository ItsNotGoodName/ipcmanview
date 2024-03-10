package main

import (
	"context"
	"errors"
	"fmt"
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
}

var mainCmd struct {
	LoggingLevel string `env:"LOGGING_LEVEL" enum:"debug,info,warn,error" default:"info"`
	LoggingType  string `env:"LOGGING_TYPE" enum:"json,console" default:"console"`

	Debug_  CmdDebug   `name:"debug" cmd:""`
	Serve   CmdServe   `cmd:"" help:"Start application." default:"1"`
	Version CmdVersion `cmd:"" help:"Show version."`
}

func main() {
	godotenv.Load()

	ktx := kong.Parse(&mainCmd, kong.Description("Application for managing and viewing Dahua devices."))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	initLogger(mainCmd.LoggingLevel, mainCmd.LoggingType)

	err := ktx.Run(&Context{
		Context: ctx,
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal().Err(err).Send()
	}
}

func initLogger(level string, typ string) {
	// Get levels
	var zerologLevel zerolog.Level
	var slogLevel slog.Level
	switch level {
	case "debug":
		zerologLevel = zerolog.DebugLevel
		slogLevel = slog.LevelDebug
	case "info":
		zerologLevel = zerolog.InfoLevel
		slogLevel = slog.LevelInfo
	case "warn":
		zerologLevel = zerolog.WarnLevel
		slogLevel = slog.LevelWarn
	case "error":
		zerologLevel = zerolog.ErrorLevel
		slogLevel = slog.LevelError
	default:
		panic(fmt.Sprintf("invalid logging level: '%s'", level))
	}

	// Set logger
	switch typ {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerologLevel)
	case "json":
		log.Logger = log.Level(zerologLevel)
	default:
		panic(fmt.Sprintf("invalid logging type: '%s'", level))
	}
	slog.SetDefault(slog.New(slogzerolog.Option{Level: slogLevel, Logger: &log.Logger}.NewZerologHandler()))
}
