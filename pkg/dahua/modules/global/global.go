package global

import (
	"github.com/ItsNotGoodName/pkg/dahua"
)

func FirstLogin(gen dahua.Generator, username string) (dahua.Response[dahua.AuthParam], error) {
	return dahua.SendRaw[dahua.AuthParam](gen.
		RPCLogin().
		Method("global.login").
		Params(dahua.JSON{
			"userName":   username,
			"password":   "",
			"loginType":  "Direct",
			"clientType": "Web3.0",
		}))
}

func SecondLogin(gen dahua.Generator, username, password, loginType, authorityType string) error {
	_, err := dahua.Send[any](gen.
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

func GetCurrentTime(gen dahua.Generator) (string, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return "", err
	}

	res, err := dahua.Send[struct {
		Time string `json:"time"`
	}](rpc.Method("global.getCurrentTime"))

	return res.Params.Time, err
}

func KeepAlive(gen dahua.Generator) (int, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return 0, err
	}

	res, err := dahua.Send[struct {
		Timeout int `json:"timeout"`
	}](rpc.Method("global.keepAlive"))

	return res.Params.Timeout, err
}

func Logout(gen dahua.Generator) (bool, error) {
	rpc, err := gen.RPC()
	if err != nil {
		return false, err
	}

	res, err := dahua.Send[bool](rpc.Method("global.logout"))

	return res.Params, err
}
