package main

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EnvUsername = "IPC_USERNAME"
	EnvPassword = "IPC_PASSWORD"
	EnvIPIn     = "IPC_IP_IN"
	EnvIPOut    = "IPC_IP_OUT"
)

func main() {
	ctx, cancel := interrupt.Context()
	defer cancel()

	inIp := getEnv(EnvIPIn)
	outIp := getEnv(EnvIPOut)
	username := getEnv(EnvUsername)
	password := getEnv(EnvPassword)

	fmt.Printf("Testing audio features from %s to %s\n", inIp, outIp)

	inClient := dahuacgi.NewConn(http.Client{}, inIp, username, password)
	clientOut := dahuacgi.NewConn(http.Client{}, outIp, username, password)

	inCount, err := dahuacgi.AudioInputChannelCount(ctx, inClient)
	if err != nil {
		log.Err(err).Msgf("Failed to get audio input count for %s", inIp)
	} else {
		fmt.Printf("Audio input count for %s = %d\n", inIp, inCount)
	}

	outCount, err := dahuacgi.AudioOutputChannelCount(ctx, clientOut)
	if err != nil {
		log.Err(err).Msgf("Failed to get audio output count for %s", outIp)
	} else {
		fmt.Printf("Audio output count for %s = %d\n", outIp, outCount)
	}

	if inCount == 0 || outCount == 0 {
		fmt.Println("Skipping audio test because of unsupported audio input/output count")
	} else {
		// Stream from cgiIn -> cgiOut
		fmt.Println("Getting audio")
		stream, err := dahuacgi.AudioStreamGet(ctx, inClient, 0, dahuacgi.HTTPTypeSinglePart)
		if err != nil {
			log.Err(err).Msgf("Failed to get audio stream for %s", inIp)
		} else {
			_, wt := io.Pipe()

			go func() {
				fmt.Println("Posting audio")
				panic("not implemented")
				// err = dahuacgi.AudioStreamPost(ctx, clientOut, 0, dahuacgi.HTTPTypeSinglePart, stream.ContentType, rd)
				if err != nil {
					log.Err(err).Msgf("Failed to post audio stream to %s", outIp)
				}
			}()

			for i := 0; true; i-- {
				copied, err := io.Copy(wt, stream)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to copy from input")
				}
				fmt.Println("Finished with written", copied)
			}
			wt.Close()
			time.Sleep(1 * time.Second)
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
