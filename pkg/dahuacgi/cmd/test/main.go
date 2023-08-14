package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvUsername = "IPC_USERNAME"
	EnvPassword = "IPC_PASSWORD"
	EnvIP       = "IPC_IP"
)

func main() {
	ctx, cancel := interrupt.Context()
	defer cancel()

	ip := getEnv(EnvIP)
	username := getEnv(EnvUsername)
	password := getEnv(EnvPassword)

	fmt.Printf("Testing CGI features on %s\n", ip)

	cgi := dahuacgi.NewConn(ip, username, password)

	// Audio
	fmt.Println("Testing audio...")
	if count, err := dahuacgi.AudioInputChannelCount(ctx, cgi); err != nil {
		log.Err(err).Msg("Failed to list audio input")
	} else {
		fmt.Println("Audio input count", count)
	}
	if count, err := dahuacgi.AudioOutputChannelCount(ctx, cgi); err != nil {
		log.Err(err).Msg("Failed to list audio output")
	} else {
		fmt.Println("Audio output count", count)
	}
	if stream, err := dahuacgi.AudioStreamGet(ctx, cgi, 0, dahuacgi.HTTPTypeSinglePart); err != nil {
		log.Err(err).Msg("Failed to get audio output")
	} else {
		b := make([]byte, 1024)

		fmt.Println("Audio stream is", stream.ContentType)
		fmt.Println("Audio streaming...")
		for i := 0; i < 10; i++ {
			audio, err := stream.Read(b)
			if err != nil {
				log.Err(err).Msg("Failed to read audio stream")
				break
			}
			fmt.Println("Read", audio, "bytes")
		}
		stream.Close()
	}

	// Snapshot
	fmt.Println("Testing snapshot...")
	if snapshot, err := dahuacgi.SnapshotGet(ctx, cgi, 0); err != nil {
		log.Err(err).Msg("Failed to snapshot")
	} else {
		defer snapshot.Close()
		written, err := io.Copy(io.Discard, snapshot)
		if err != nil {
			log.Err(err).Msg("Failed to read snapshot")
		}
		snapshot.Close()

		fmt.Printf("Snapshot is %d bytes\n", written)
	}

	// Events
	fmt.Println("Testing events forever...")
	if eventManager, err := dahuacgi.EventManagerGet(ctx, cgi, 0); err != nil {
		log.Err(err).Msg("Failed to events")
	} else {
		eventSession := eventManager.Reader()
		for {
			fmt.Println("Waiting for next event...")

			err := eventSession.Poll()
			if err != nil {
				log.Err(err).Msg("Failed to SeekBoundary")
				break
			}

			event, err := eventSession.ReadEvent()
			if err != nil {
				log.Err(err).Msg("Failed to parse event")
			}
			data := string(event.Data)

			fmt.Printf("Event: %+v\n", event)
			fmt.Printf("Event Data: %s\n", data)
		}
	}
}

func getEnv(env string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		log.Fatal().Str("env", env).Msg("Environment variable not set")
	}

	return str
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
