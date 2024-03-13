package gorise

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func BuildNtfy(cfg Config) (Sender, error) {
	return buildNtfy(false, cfg)
}

func BuildNtfys(cfg Config) (Sender, error) {
	return buildNtfy(true, cfg)
}

func buildNtfy(https bool, cfg Config) (Sender, error) {
	paths := cfg.Paths()
	pathsLen := len(paths)
	if pathsLen == 0 {
		return nil, fmt.Errorf("ntfy: no config")
	}
	if pathsLen > 2 {
		return nil, fmt.Errorf("ntfy: multiple topics are not supported")
	}

	protocol := "http://"
	if https {
		protocol = "https://"
	}

	var (
		host  string
		topic string
	)
	if pathsLen == 2 {
		host = paths[0]
		topic = paths[1]
	} else {
		topic = paths[0]
	}

	return NewNtfy(protocol+host, topic), nil
}

func NewNtfy(url, topic string) Ntfy {
	if url == "" {
		url = "https://ntfy.sh"
	}
	return Ntfy{
		url: fmt.Sprintf("%s/%s", url, topic),
	}
}

type Ntfy struct {
	// authorization string
	url string
}

func (n Ntfy) Send(ctx context.Context, msg Message) error {
	text := msg.Text()
	if text != "" {
		req, err := http.NewRequestWithContext(ctx, "POST", n.url, strings.NewReader(text))
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(res.Body)
			res.Body.Close()
			return errors.New(string(b))
		}
		res.Body.Close()
	}

	for _, a := range msg.Attachments {
		err := func() error {

			req, err := http.NewRequestWithContext(ctx, "PUT", n.url, a.Reader)
			if err != nil {
				return err
			}

			req.Header.Set("Filename", a.Name)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(res.Body)
				return errors.New(string(b))
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
