package dahuarpc

import (
	"context"
	"encoding/json"
)

type Instance struct {
	method    string
	responses map[string]Response[json.RawMessage]
}

func NewInstance(method string) *Instance {
	return &Instance{
		method:    method,
		responses: make(map[string]Response[json.RawMessage]),
	}
}

func (i *Instance) Get(ctx context.Context, c Client, key string, params any) (Response[json.RawMessage], error) {
	if res, found := i.responses[key]; found {
		return res, nil
	}

	rpc, err := c.RPC(ctx)
	if err != nil {
		return Response[json.RawMessage]{}, err
	}

	res, err := Send[json.RawMessage](ctx, rpc.Method(i.method).Params(params))
	if err != nil {
		return Response[json.RawMessage]{}, err
	}

	i.responses[key] = res

	return res, nil
}
