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

func SnapshotGet(ctx context.Context, c Conn, channel int) (Snapshot, error) {
	req := NewRequest("snapshot.cgi")

	if channel != 0 {
		req.QueryInt("channel", channel)
	}

	res, err := OK(c.CGIGet(ctx, req))
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
