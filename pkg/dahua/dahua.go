package dahua

import (
	"bytes"
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

type Generator interface {
	RPC() (RequestBuilder, error)
	RPCLogin() RequestBuilder
}

type ResponseSession struct {
	Value string
}

func (s *ResponseSession) UnmarshalJSON(data []byte) error {
	// string -> string
	if err := json.Unmarshal(data, &s.Value); err == nil {
		return nil
	}

	// int64 -> string
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		s.Value = strconv.FormatInt(num, 10)
		return nil
	}

	return fmt.Errorf("session is not a string or number")
}

type ResponseResult struct {
	Number int64
}

func (s *ResponseResult) UnmarshalJSON(data []byte) error {
	// int64 -> int64
	if err := json.Unmarshal(data, &s.Number); err == nil {
		return nil
	}

	// bool -> int64
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			s.Number = 1
		}
		return nil
	}

	return fmt.Errorf("result is not a number or boolean")
}

func (s *ResponseResult) Bool() bool {
	return s.Number == 1
}

// Response from the camera.
// T should always be as pointer type, this will prevent null being defaulted to T.
type Response[T any] struct {
	ID      int             `json:"id"`
	Session ResponseSession `json:"session"` // Session can be a string or an int64
	Error   *ResponseError  `json:"error"`
	Params  T               `json:"params"`
	Result  ResponseResult  `json:"result"` // Result can be a bool or an int64
}

type ResponseError struct {
	Code    int
	Message string
	Kind    ErrResponseKind
}

func (r *ResponseError) UnmarshalJSON(data []byte) error {
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
		r.Kind = ErrResponseKindInvalidRequest
	case 268894210:
		r.Kind = ErrResponseMethodNotFound
	case 268632064:
		r.Kind = ErrResponseInterfaceNotFound
	case 285409284:
		r.Kind = ErrResponseNoData
	default:
		r.Kind = ErrResponseKindUnknown
	}

	return nil
}

func (r *ResponseError) Error() string {
	return r.Message
}

type ErrResponseKind = string

var (
	ErrResponseKindInvalidRequest ErrResponseKind = "InvalidRequest"
	ErrResponseMethodNotFound     ErrResponseKind = "MethodNotFound"
	ErrResponseInterfaceNotFound  ErrResponseKind = "InterfaceNotFound"
	ErrResponseNoData             ErrResponseKind = "NoData"
	ErrResponseKindUnknown        ErrResponseKind = "Unknown"
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

// TODO: I want this attached to the RequestBuilder
func SendRaw[T any](r RequestBuilder) (Response[T], error) {
	var res Response[T]

	b, err := json.Marshal(r.req)
	if err != nil {
		return res, err
	}

	resp, err := r.client.Post(r.url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return res, errors.Join(ErrRequestFailed, err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, err
	}

	return res, nil
}

// TODO: I want this attached to the RequestBuilder
func Send[T any](r RequestBuilder) (Response[T], error) {
	res, err := SendRaw[T](r)
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
