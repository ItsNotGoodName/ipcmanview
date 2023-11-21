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

func NewConn(httpClient http.Client, httpAddress, username, password string) Conn {
	t := &digest.Transport{
		Username: username,
		Password: password,
	}
	if httpClient.Transport != nil {
		t.Transport = httpClient.Transport
	}
	httpClient.Transport = t
	return Conn{
		baseURL: fmt.Sprintf("%s/cgi-bin/", httpAddress),
		client:  &httpClient,
	}
}

func (c Conn) CGIGet(ctx context.Context, r *Request) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.URL(c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(r.Request(req))
}
