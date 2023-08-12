package dahuacgi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/icholy/digest"
)

var _ GenCGI = (*Conn)(nil)

type Conn struct {
	client  *http.Client
	baseURL string
}

func NewConn(ip, username, password string) Conn {
	return Conn{
		baseURL: fmt.Sprintf("http://%s/cgi-bin/", ip),
		client: &http.Client{
			Transport: &digest.Transport{
				Username: username,
				Password: password,
			},
		},
	}
}

func (c Conn) CGIGet(ctx context.Context, method string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+method, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c Conn) CGIPost(ctx context.Context, method string, headers http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+method, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers

	return c.client.Do(req)
}
