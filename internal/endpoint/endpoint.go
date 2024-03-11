package endpoint

import (
	"context"
	"fmt"
	"io"
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

var Builders = map[string]func(urL string) (Sender, error){
	"tgram": BuildTelegram,
}

func Build(urL string) (Sender, error) {
	scheme, _, _ := strings.Cut(urL, "://")
	builder, ok := Builders[scheme]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrBuilderNotFound, scheme)
	}
	return builder(urL)
}
