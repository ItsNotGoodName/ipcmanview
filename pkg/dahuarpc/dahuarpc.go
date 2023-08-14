// dahuarpc is a RPC client library for Dahua's RPC API.
package dahuarpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	ErrInvalidSession = fmt.Errorf("invalid session")
	ErrRequestFailed  = fmt.Errorf("request failed")
)

type Client interface {
	RPC(ctx context.Context) (RequestBuilder, error)
}

type ClientLogin interface {
	RPCLogin() RequestBuilder
}

type ResponseSession string

func (s *ResponseSession) UnmarshalJSON(data []byte) error {
	// string -> string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = ResponseSession(str)
		return nil
	}

	// int64 -> string
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*s = ResponseSession(strconv.FormatInt(num, 10))
		return nil
	}

	return fmt.Errorf("session is not a string or number")
}

func (s ResponseSession) String() string {
	return string(s)
}

type ResponseResult int64

func (s *ResponseResult) UnmarshalJSON(data []byte) error {
	// int64 -> int64
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		*s = ResponseResult(num)
		return nil
	}

	// bool -> int64
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*s = 1
		}
		return nil
	}

	return fmt.Errorf("result is not a number or boolean")
}

func (s ResponseResult) Integer() int64 {
	return int64(s)
}

func (s ResponseResult) Bool() bool {
	return s == 1
}

// Response from the camera.
type Response[T any] struct {
	ID      int             `json:"id"`
	Session ResponseSession `json:"session"`
	Error   *ErrResponse    `json:"error"`
	Params  T               `json:"params"`
	Result  ResponseResult  `json:"result"`
}

type ErrResponse struct {
	Code    int
	Message string
	Type    ErrResponseType
}

func (r *ErrResponse) UnmarshalJSON(data []byte) error {
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	r.Code = res.Code
	r.Message = res.Message

	switch res.Code {
	case 268894209:
		r.Type = ErrResponseTypeInvalidRequest
	case 268894210:
		r.Type = ErrResponseTypeMethodNotFound
	case 268632064:
		r.Type = ErrResponseTypeInterfaceNotFound
	case 285409284:
		r.Type = ErrResponseTypeNoData
	default:
		r.Type = ErrResponseTypeUnknown
	}

	return nil
}

func (r *ErrResponse) Error() string {
	return r.Message
}

type ErrResponseType string

var (
	ErrResponseTypeInvalidRequest    ErrResponseType = "InvalidRequest"
	ErrResponseTypeMethodNotFound    ErrResponseType = "MethodNotFound"
	ErrResponseTypeInterfaceNotFound ErrResponseType = "InterfaceNotFound"
	ErrResponseTypeNoData            ErrResponseType = "NoData"
	ErrResponseTypeUnknown           ErrResponseType = "Unknown"
)

type Request struct {
	ID      int    `json:"id"`
	Session string `json:"session,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	Object  int64  `json:"object,omitempty"`
}

type RequestBuilder struct {
	client *http.Client
	req    Request
	url    string
}

func NewRequestBuilder(client *http.Client, id int, url, session string) RequestBuilder {
	return RequestBuilder{
		client: client,
		req: Request{
			ID:      id,
			Session: session,
		},
		url: url,
	}
}

func (r RequestBuilder) Params(params any) RequestBuilder {
	r.req.Params = params
	return r
}

func (r RequestBuilder) Object(object int64) RequestBuilder {
	r.req.Object = object
	return r
}

func (r RequestBuilder) Method(method string) RequestBuilder {
	r.req.Method = method
	return r
}

// SendRaw sends RPC request to camera without checking if the response contains an error field.
func SendRaw[T any](ctx context.Context, r RequestBuilder) (Response[T], error) {
	var res Response[T]

	b, err := json.Marshal(r.req)
	if err != nil {
		return res, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(b))
	if err != nil {
		return res, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return res, errors.Join(ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, err
	}

	return res, nil
}

// Send RPC request to camera and check the response's error field.
func Send[T any](ctx context.Context, r RequestBuilder) (Response[T], error) {
	res, err := SendRaw[T](ctx, r)
	if err != nil {
		return res, err
	}
	if res.Error != nil {
		if res.Error.Code == 287637505 || res.Error.Code == 287637504 {
			return res, errors.Join(ErrInvalidSession, res.Error)
		}
		return res, res.Error
	}

	return res, nil
}
