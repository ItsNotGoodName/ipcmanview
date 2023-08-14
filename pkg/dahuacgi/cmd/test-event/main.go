package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed fixture.txt
var fixture string

func main() {
	eventReader := dahuacgi.NewEventReader(strings.NewReader(fixture), dahuacgi.DefaultEventBoundary)

	for i := 0; ; i++ {
		err := eventReader.Poll()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal().Err(err).Msg("Failed to seek next SeekBoundary")
		}

		event, err := eventReader.ReadEvent()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse event")
		}

		fmt.Printf("%+v\n", event)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
