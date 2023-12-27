// dahuarpc is a client library for Dahua's RPC API.
package dahuarpc

import (
	"context"
	"encoding/json"
	"io"
)

type Conn interface {
	Do(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error)
	// SessionRaw() string
}

// SendRaw sends RPC request to camera without checking if the response contains an error field.
func SendRaw[T any](ctx context.Context, c Conn, rb RequestBuilder) (Response[T], error) {
	var res Response[T]

	rd, err := c.Do(ctx, rb)
	if err != nil {
		return res, err
	}
	defer rd.Close()

	b, err := io.ReadAll(rd)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(b, &res); err != nil {
		return res, err
	}

	if res.Error != nil {
		res.Error.Method = rb.Request.Method
	}

	return res, nil
}

// Send RPC request to device and check the response's error field.
func Send[T any](ctx context.Context, c Conn, rb RequestBuilder) (Response[T], error) {
	res, err := SendRaw[T](ctx, c, rb)
	if err != nil {
		return res, err
	}
	if res.Error != nil {
		return res, res.Error
	}

	return res, nil
}
