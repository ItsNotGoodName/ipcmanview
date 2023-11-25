package global

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Client interface {
	RawRPC(ctx context.Context) (dahuarpc.RequestBuilder, error)
	RawRPCLogin() dahuarpc.RequestBuilder
}

func FirstLogin(ctx context.Context, c Client, username string) (dahuarpc.Response[dahuarpc.AuthParam], error) {
	return dahuarpc.SendRaw[dahuarpc.AuthParam](ctx, c.
		RawRPCLogin().
		Method("global.login").
		Params(struct {
			Username   string `json:"userName"`
			Password   string `json:"password"`
			LoginType  string `json:"loginType"`
			ClientType string `json:"clientType"`
		}{
			Username:   username,
			Password:   "",
			LoginType:  "Direct",
			ClientType: "Web3.0",
		}))
}

func SecondLogin(ctx context.Context, c Client, username, password, loginType, authorityType string) error {
	_, err := dahuarpc.Send[any](ctx, c.
		RawRPCLogin().
		Method("global.login").
		Params(struct {
			Username      string `json:"userName"`
			Password      string `json:"password"`
			LoginType     string `json:"loginType"`
			ClientType    string `json:"clientType"`
			AuthorityType string `json:"authorityType"`
		}{
			Username:      username,
			Password:      password,
			LoginType:     loginType,
			ClientType:    "Web3.0",
			AuthorityType: authorityType,
		}))

	return err
}

func GetCurrentTime(ctx context.Context, c Client) (string, error) {
	rpc, err := c.RawRPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahuarpc.Send[struct {
		Time string `json:"time"`
	}](ctx, rpc.Method("global.getCurrentTime"))

	return res.Params.Time, err
}

func KeepAlive(ctx context.Context, c Client) (int, error) {
	rpc, err := c.RawRPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahuarpc.Send[struct {
		Timeout int `json:"timeout"`
	}](ctx, rpc.Method("global.keepAlive"))

	return res.Params.Timeout, err
}

func Logout(ctx context.Context, c Client) (bool, error) {
	rpc, err := c.RawRPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahuarpc.Send[bool](ctx, rpc.Method("global.logout"))

	return res.Params, err
}
