package endpoint

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
)

var (
	ErrBuilderNotFound = fmt.Errorf("builder not found")
)

type Message struct {
	Title       string
	Body        string
	Attachments []Attachment
}

type Attachment struct {
	Name string
	Mime string
	io.Reader
}

func (a Attachment) IsImage() bool {
	return strings.HasPrefix(a.Mime, "image/")
}

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

var Builders = map[string]func(urL *url.URL) (Sender, error){
	"tgram": TelegramFromURL,
}

func SenderFromURL(urL *url.URL) (Sender, error) {
	builder, ok := Builders[urL.Scheme]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrBuilderNotFound, urL.Scheme)
	}
	return builder(urL)
}
