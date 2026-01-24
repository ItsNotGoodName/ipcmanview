package dahuarpc

import "encoding/json"

type Request struct {
	ID      int
	Session string
	Method  string
	Params  any
	Object  int64
	// TODO: don't use map
	options map[string]any
}

func (r Request) MarshalJSON() ([]byte, error) {
	// ID
	r.options["id"] = r.ID

	// Session
	if r.Session != "" {
		r.options["session"] = r.Session
	} else {
		delete(r.options, "session")
	}

	// Method
	r.options["method"] = r.Method

	// Params
	params, err := json.Marshal(r.Params)
	if err != nil {
		return nil, err
	}
	r.options["params"] = json.RawMessage(params)

	// Object
	if r.Object != 0 {
		r.options["object"] = r.Object
	} else {
		delete(r.options, "object")
	}

	return json.Marshal(r.options)
}

func New(method string) RequestBuilder {
	return RequestBuilder{
		Login: false,
		Request: Request{
			ID:      0,
			Session: "",
			Method:  method,
			Params:  nil,
			options: make(map[string]any),
		},
	}
}

func NewLogin(method string) RequestBuilder {
	rb := New(method)
	rb.Login = true
	return rb
}

type RequestBuilder struct {
	Login   bool
	Request Request
}

func (rb RequestBuilder) ID(id int) RequestBuilder {
	rb.Request.ID = id
	return rb
}

func (rb RequestBuilder) Session(session string) RequestBuilder {
	rb.Request.Session = session
	return rb
}

func (rb RequestBuilder) Params(params any) RequestBuilder {
	rb.Request.Params = params
	return rb
}

func (rb RequestBuilder) Object(object int64) RequestBuilder {
	rb.Request.Object = object
	return rb
}

func (rb RequestBuilder) Option(key string, value any) RequestBuilder {
	rb.Request.options[key] = value
	return rb
}
