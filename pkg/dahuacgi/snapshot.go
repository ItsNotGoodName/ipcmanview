package dahuacgi

import (
	"context"
	"io"
)

type Snapshot struct {
	io.ReadCloser
	ContentType   string
	ContentLength string
}

func SnapshotGet(ctx context.Context, c Conn, channel, typE int) (Snapshot, error) {
	req := New("snapshot.cgi")

	if channel != 0 {
		req.QueryInt("channel", channel)
	}
	if typE != 0 {
		req.QueryInt("type", typE)
	}

	res, err := OK(c.Do(ctx, req))
	if err != nil {
		return Snapshot{}, err
	}

	contentType := res.Header.Get("Content-Type")
	contentLength := res.Header.Get("Content-Length")

	return Snapshot{
		ReadCloser:    res.Body,
		ContentType:   contentType,
		ContentLength: contentLength,
	}, nil
}
