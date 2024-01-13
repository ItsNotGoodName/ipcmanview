package dahuarpc

import (
	"context"
	"encoding/json"
)

func NewCache() Cache {
	return Cache(make(map[string]Response[json.RawMessage]))
}

// Cache caches RPC calls.
type Cache map[string]Response[json.RawMessage]

func (c Cache) Send(ctx context.Context, conn Conn, key string, rb RequestBuilder) (Response[json.RawMessage], error) {
	if res, found := c[key]; found {
		return res, nil
	}

	res, err := Send[json.RawMessage](ctx, conn, rb)
	if err != nil {
		return Response[json.RawMessage]{}, err
	}

	c[key] = res

	return res, nil
}
