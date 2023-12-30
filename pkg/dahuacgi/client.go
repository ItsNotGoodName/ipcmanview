package dahuacgi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/icholy/digest"
)

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(httpClient http.Client, u *url.URL, username, password string) Client {
	t := &digest.Transport{
		Username: username,
		Password: password,
	}
	if httpClient.Transport != nil {
		t.Transport = httpClient.Transport
	}
	httpClient.Transport = t
	return Client{
		baseURL: fmt.Sprintf("%s://%s/cgi-bin/", u.Scheme, u.Hostname()),
		client:  &httpClient,
	}
}

func (c Client) Do(ctx context.Context, r *Request) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.URL(c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(r.Request(req))
}
