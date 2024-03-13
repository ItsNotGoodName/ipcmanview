package gorise

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"os/signal"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/attachment.jpg
var attachment []byte

func message() Message {
	return Message{
		Title: "Test title",
		Body:  "Test body.",
		Attachments: []Attachment{
			{
				Name:   "test.jpg",
				Mime:   "image/jpeg",
				Reader: bytes.NewBuffer(attachment),
			},
		},
	}
}

func TestGorise(t *testing.T) {
	var tests []string = []string{
		"GORISE_TEST_TELEGRAM_URL",
		"GORISE_TEST_NTFY_URL",
	}

	godotenv.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	for _, tt := range tests {
		url, ok := os.LookupEnv(tt)
		if !ok {
			continue
		}

		t.Run(tt, func(t *testing.T) {
			sender, err := Build(url)
			if !assert.NoError(t, err) {
				return
			}

			err = sender.Send(ctx, message())
			if !assert.NoError(t, err) {
				return
			}
		})
	}
}
