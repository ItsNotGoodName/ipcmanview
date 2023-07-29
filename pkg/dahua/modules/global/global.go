package global

import (
	"context"

	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
)

func FirstLogin(ctx context.Context, gen dahua.GenRPCLogin, username string) (dahua.Response[dahua.AuthParam], error) {
	return dahua.SendRaw[dahua.AuthParam](ctx, gen.
		RPCLogin().
		Method("global.login").
		Params(dahua.JSON{
			"userName":   username,
			"password":   "",
			"loginType":  "Direct",
			"clientType": "Web3.0",
		}))
}

func SecondLogin(ctx context.Context, gen dahua.GenRPCLogin, username, password, loginType, authorityType string) error {
	_, err := dahua.Send[any](ctx, gen.
		RPCLogin().
		Method("global.login").
		Params(dahua.JSON{
			"userName":      username,
			"password":      password,
			"clientType":    "Web3.0",
			"loginType":     loginType,
			"authorityType": authorityType,
		}))

	return err
}

func GetCurrentTime(ctx context.Context, gen dahua.GenRPC) (string, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Time string `json:"time"`
	}](ctx, rpc.Method("global.getCurrentTime"))

	return res.Params.Time, err
}

func KeepAlive(ctx context.Context, gen dahua.GenRPC) (int, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return 0, err
	}

	res, err := dahua.Send[struct {
		Timeout int `json:"timeout"`
	}](ctx, rpc.Method("global.keepAlive"))

	return res.Params.Timeout, err
}

func Logout(ctx context.Context, gen dahua.GenRPC) (bool, error) {
	rpc, err := gen.RPC(ctx)
	if err != nil {
		return false, err
	}

	res, err := dahua.Send[bool](ctx, rpc.Method("global.logout"))

	return res.Params, err
}
