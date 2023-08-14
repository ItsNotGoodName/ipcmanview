package dahuacgi

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client interface {
	CGIGet(ctx context.Context, req *Request) (*http.Response, error)
}

type Request struct {
	method string
	query  url.Values
	Header http.Header
}

func NewRequest(method string) *Request {
	return &Request{
		method: method,
		query:  url.Values{},
		Header: http.Header{},
	}
}

func (r *Request) QueryString(key string, value string) *Request {
	r.query.Add(key, value)
	return r
}

func (r *Request) QueryInt(key string, value int) *Request {
	r.query.Add(key, strconv.Itoa(value))
	return r
}

func (r *Request) HeaderString(key string, value string) *Request {
	r.Header.Add(key, value)
	return r
}

func (r *Request) URL(baseURL string) string {
	query := r.query.Encode()
	if query != "" {
		query = "?" + query
	}
	return baseURL + r.method + query
}

func (r *Request) Request(req *http.Request) *http.Request {
	req.Header = r.Header
	return req
}

func OK(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}

	// OK
	if res.StatusCode < 200 || res.StatusCode > 299 {
		res.Body.Close()
		return nil, fmt.Errorf(res.Status)
	}

	return res, nil
}

type Table []TableData

type TableData struct {
	Key   string
	Value string
}

func OKTable(res *http.Response, err error) (Table, error) {
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// OK
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf(res.Status)
	}

	// Table
	sc := bufio.NewScanner(res.Body)
	var table Table
	for sc.Scan() {
		kv := strings.SplitN(sc.Text(), "=", 2)
		if len(kv) != 2 {
			continue
		}
		table = append(table, TableData{Key: kv[0], Value: kv[1]})
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return table, nil
}

func (t Table) Get(key string) string {
	for i := range t {
		if t[i].Key == key {
			return t[i].Value
		}
	}

	return ""
}
