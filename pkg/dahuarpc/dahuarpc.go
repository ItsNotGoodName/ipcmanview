// dahuarpc is a client library for Dahua's RPC API.
package dahuarpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrRequestFailed = errors.New("request failed")
)

const (
	StateLogout State = iota
	StateLogin
	StateError
	StateClosed
)

type State int

func (s State) String() string {
	switch s {
	case StateLogin:
		return "login"
	case StateLogout:
		return "logout"
	case StateError:
		return "error"
	case StateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

type Conn interface {
	Do(ctx context.Context, rb RequestBuilder) (io.ReadCloser, error)
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

	// fmt.Printf("RESPONSE: %s\n", string(b))

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

// DoRaw executes the RPC request and returns the body.
func DoRaw(ctx context.Context, rb RequestBuilder, httpClient *http.Client, urL string) (io.ReadCloser, error) {
	b, err := json.Marshal(rb.Request)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("REQUEST: %s\n", string(b))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Join(ErrRequestFailed, err)
	}
	return resp.Body, nil
}
