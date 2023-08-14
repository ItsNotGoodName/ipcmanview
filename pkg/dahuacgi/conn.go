package dahuacgi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/icholy/digest"
)

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

func (c Conn) CGIGet(ctx context.Context, r *Request) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.URL(c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(r.Request(req))
}
