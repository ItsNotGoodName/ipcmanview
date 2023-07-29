package global

import "github.com/ItsNotGoodName/pkg/dahua"

func FirstLogin(gen dahua.Generator, username string) (dahua.Response[dahua.AuthParam], error) {
	a := gen.
		RPCLogin().
		Method("global.login").
		Params(dahua.JSON{
			"userName":   username,
			"password":   "",
			"loginType":  "Direct",
			"clientType": "Web3.0",
		})

	return dahua.SendRaw[dahua.AuthParam](a)
}

func SecondLogin(gen dahua.Generator, username, password, loginType, authorityType string) error {
	a := gen.
		RPCLogin().
		Method("global.login").
		Params(dahua.JSON{
			"userName":      username,
			"password":      password,
			"clientType":    "Web3.0",
			"loginType":     loginType,
			"authorityType": authorityType,
		})

	_, err := dahua.Send[any](a)

	return err
}

func GetCurrentTime(gen dahua.Generator) (string, error) {
	a := gen.
		RPC().
		Method("global.getCurrentTime")

	b, err := dahua.Send[struct {
		Time string `json:"time"`
	}](a)

	return b.Params.Time, err
}

func KeepAlive(gen dahua.Generator) (int, error) {
	a := gen.
		RPC().
		Method("global.keepAlive")

	b, err := dahua.Send[struct {
		Timeout int `json:"timeout"`
	}](a)

	return b.Params.Timeout, err
}

func Logout(gen dahua.Generator) (bool, error) {
	a := gen.
		RPC().
		Method("global.logout")

	b, err := dahua.Send[bool](a)

	return b.Params, err
}
