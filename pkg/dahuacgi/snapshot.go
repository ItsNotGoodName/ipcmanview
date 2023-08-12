package dahuacgi

import (
	"context"
	"io"
	"net/url"
	"strconv"
)

type Snapshot struct {
	io.ReadCloser
	ContentType   string
	ContentLength string
}

func SnapshotGet(ctx context.Context, cgi GenCGI, channel int) (Snapshot, error) {
	method := "snapshot.cgi"

	query := url.Values{}
	query.Add("action", "attach")
	query.Add("codes", "[All]")
	if channel != 0 {
		query.Add("channel", strconv.Itoa(channel))
	}
	if len(query) > 0 {
		method += "?" + query.Encode()
	}
	res, err := OK(cgi.CGIGet(ctx, method))
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
