package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Flags
	ip := flag.String("ip", "", "IP address of Dahua IP camera.")
	username := flag.String("username", "admin", "Username for Dahua IP camera.")
	password := flag.String("password", "", "Password for Dahua IP camera.")
	pretty := flag.Bool("pretty", false, "Pretty print JSON.")
	timeout := flag.Int("timeout", 10, "HTTP timeout.")
	restart := flag.Bool("restart", false, "Restart when EOF error occurs.")

	// Parse
	flag.Parse()
	if *ip == "" {
		log.Fatalln("IP address not supplied.")
	}
	cfg := struct {
		IP       string
		Username string
		Password string
		Pretty   bool
		Timeout  int
		Restart  bool
	}{
		IP:       *ip,
		Username: *username,
		Password: *password,
		Pretty:   *pretty,
		Timeout:  *timeout,
		Restart:  *restart,
	}
	marshal := (func() func(v any) ([]byte, error) {
		if cfg.Pretty {
			return func(v any) ([]byte, error) {
				return json.MarshalIndent(v, "", "  ")
			}
		}

		return json.Marshal
	})()

	// Connect
	httpClient := http.Client{}
	conn := dahuacgi.NewConn(httpClient, cfg.IP, cfg.Username, cfg.Password)
	for {
		slog.Info("Connecting to Dahua IP camera", "ip", cfg.IP, "username", cfg.Username)
		stream, err := dahuacgi.EventManagerGet(ctx, conn, 0)
		if err != nil {
			log.Fatalln("Failed to get event manager:", err)
		}
		reader := stream.Reader()
		slog.Info("Connected")

		// Read events
		for id := 0; ; id++ {
			err := reader.Poll()
			if err != nil {
				if errors.Is(err, io.EOF) && cfg.Restart {
					slog.Info("Restarting", "err", err)
					break
				}
				log.Fatalln("Failed to poll reader:", err)
			}

			event, err := reader.ReadEvent()
			if err != nil {
				if errors.Is(err, io.EOF) && cfg.Restart {
					slog.Info("Restarting", "err", err)
					break
				}
				log.Fatalln("Failed to read event:", err)
			}

			b, err := marshal(struct {
				ID            int
				ContentType   string
				ContentLength int
				Code          string
				Action        string
				Index         int
				Data          json.RawMessage
			}{
				ID:            id,
				ContentType:   event.ContentType,
				ContentLength: event.ContentLength,
				Code:          event.Code,
				Action:        event.Action,
				Index:         event.Index,
				Data:          event.Data,
			})
			if err != nil {
				log.Fatalln("Failed to marshal event:", err)
			}

			fmt.Println(string(b))
		}
	}
}
