package dahuarpc

import (
	"context"
	"io"
	"net/http"
)

func NewFileClient(client *http.Client, concurrent int) FileClient {
	ctx, cancel := context.WithCancel(context.Background())
	return FileClient{
		ctx:    ctx,
		cancel: cancel,
		client: client,
		sema:   make(chan struct{}, concurrent),
	}
}

// FileClient handles file access to prevent "Resource is limited, open video failed!" errors.
// The client limits the number of concurrent requests and makes sure the body is completely drained.
type FileClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *http.Client
	sema   chan struct{}
}

func (c FileClient) Close() {
	c.cancel()
}

func (c FileClient) Do(ctx context.Context, urL, cookie string) (io.ReadCloser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case c.sema <- struct{}{}:
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, urL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Cookie", cookie)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return fileResponse{
		ReadCloser: resp.Body,
		ctx:        c.ctx,
		sema:       c.sema,
	}, nil
}

type fileResponse struct {
	io.ReadCloser
	ctx  context.Context
	sema <-chan struct{}
}

func (r fileResponse) Close() error {
	go func() {
		// Assume the device will not send more than 1 GB
		io.CopyN(io.Discard, r, 1024*1024*1024)
		r.ReadCloser.Close()
		<-r.sema
	}()
	return nil
}
