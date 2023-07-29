package dahua

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type Generator interface {
	RPC() RequestBuilder
	RPCLogin() RequestBuilder
}

type ResponseValue struct {
	Value string
}

func (s *ResponseValue) UnmarshalJSON(data []byte) error {
	// string -> string
	if err := json.Unmarshal(data, &s.Value); err == nil {
		return nil
	}

	// int64 -> string
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		return nil
	}
	s.Value = strconv.FormatInt(num, 10)

	return nil
}

type Response[T any] struct {
	ID      int            `json:"id"`
	Session ResponseValue  `json:"session"` // Session can be a string or a int64
	Error   *ResponseError `json:"error"`
	Params  T              `json:"params"`
	Result  ResponseValue  `json:"result"` // Result can be a string or a int64
}

// func Params[T any](r Response[*T]) (T, error) {
// 	if r.Params == nil {
// 		panic("No Params")
// 	}
//
// 	return *r.Params, nil
// }
//
// func RequireParams[T any](r Response[*T]) (Response[T], error) {
// 	if r.Params == nil {
// 		panic("No Params")
// 	}
//
// 	return Response[T]{
// 		ID:      r.ID,
// 		Session: r.Session,
// 		Error:   r.Error,
// 		Params:  *r.Params,
// 		Result:  r.Result,
// 	}, nil
// }

type ResponseError struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Kind    ResponseKind `json:"kind"`
}

type ResponseKind = string

var (
	ResponseKindInvalidRequest    ResponseKind = "InvalidRequest"
	ResponseKindMethodNotFound    ResponseKind = "MethodNotFound"
	ResponseKindInterfaceNotFound ResponseKind = "InterfaceNotFound"
	ResponseKindNoData            ResponseKind = "NoData"
	ResponseKindUnknown           ResponseKind = "Unknown"
)

type Request struct {
	ID      int    `json:"id"`
	Session string `json:"session,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	Object  int64  `json:"object,omitempty"`
}

type RequestBuilder struct {
	client         *http.Client
	req            Request
	url            string
	requireSession bool
}

func NewRequestBuilder(client *http.Client, id int, url string, session string) RequestBuilder {
	return RequestBuilder{
		client: client,
		req: Request{
			ID:      id,
			Session: session,
		},
		url: url,
	}
}

func (r RequestBuilder) RequireSession() RequestBuilder {
	r.requireSession = true
	return r
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

func SendRaw[T any](r RequestBuilder) (Response[T], error) {
	if r.requireSession && r.req.Session == "" {
		panic("No Session")
	}

	var target Response[T]

	b, err := json.Marshal(r.req)
	if err != nil {
		panic("Marshal: " + err.Error())
	}

	res, err := r.client.Post(r.url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic("Post: " + err.Error())
	}

	if err := json.NewDecoder(res.Body).Decode(&target); err != nil {
		panic("Decode: " + err.Error())
	}

	return target, nil
}

func Send[T any](r RequestBuilder) (Response[T], error) {
	res, err := SendRaw[T](r)
	if err != nil {
		return res, err
	}
	if res.Error != nil {
		panic("Response Error: " + res.Error.Message)
	}

	return res, nil
}
