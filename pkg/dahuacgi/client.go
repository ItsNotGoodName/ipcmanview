package dahuacgi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/icholy/digest"
)

type Config struct {
	url string
}

type ConfigFunc func(c *Config)

func WithURL(urL string) ConfigFunc {
	return func(c *Config) {
		c.url = urL
	}
}

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(ip, username, password string, configFuncs ...ConfigFunc) Client {
	cfg := Config{
		url: fmt.Sprintf("http://%s/cgi-bin/", ip),
	}

	for _, fn := range configFuncs {
		fn(&cfg)
	}

	return Client{
		baseURL: cfg.url,
		client: &http.Client{
			Transport: &digest.Transport{
				Username: username,
				Password: password,
			},
		},
	}
}

func URL(u *url.URL) string {
	return fmt.Sprintf("%s://%s/cgi-bin/", u.Scheme, u.Hostname())
}

func (c Client) Do(ctx context.Context, r *Request) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.URL(c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(r.Request(req))
}
