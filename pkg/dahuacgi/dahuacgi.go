package dahuacgi

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Gen interface {
	CGIGet(ctx context.Context, method string) (*http.Response, error)
	CGIPost(ctx context.Context, method string, headers http.Header, body io.Reader) (*http.Response, error)
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

func OKBody(res *http.Response, err error) (io.ReadCloser, error) {
	if err != nil {
		return nil, err
	}

	// OK
	if res.StatusCode < 200 || res.StatusCode > 299 {
		res.Body.Close()
		return nil, fmt.Errorf(res.Status)
	}

	// Body
	return res.Body, nil
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
