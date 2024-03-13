package gorise

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

// Text returns the title and body combined.
func (m Message) Text() string {
	if m.Title == "" {
		return m.Body
	}
	if m.Body == "" {
		return m.Title
	}
	return m.Title + "\n" + m.Body
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

func NewURL(urL string) (URL, error) {
	scheme, data, ok := strings.Cut(urL, "://")
	if !ok {
		return URL{}, fmt.Errorf("invalid config url")
	}

	return URL{
		Scheme: scheme,
		Config: Config(data),
	}, nil
}

type URL struct {
	Scheme string
	Config Config
}

type Config string

func (c Config) Paths() []string {
	return strings.Split(string(c), "/")
}

var Builders = map[string]func(cfg Config) (Sender, error){
	"console": BuildConsole,
	"tgram":   BuildTelegram,
	"ntfy":    BuildNtfy,
	"ntfys":   BuildNtfys,
}

func Build(urL string) (Sender, error) {
	url, err := NewURL(urL)
	if err != nil {
		return nil, err
	}

	builder, ok := Builders[url.Scheme]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrBuilderNotFound, url.Scheme)
	}
	return builder(url.Config)
}
