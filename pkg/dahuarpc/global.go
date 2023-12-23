package dahuarpc

import (
	"context"
)

func FirstLogin(ctx context.Context, c Conn, username string) (Response[AuthParam], error) {
	return SendRaw[AuthParam](ctx, c, NewLogin("global.login").
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

func SecondLogin(ctx context.Context, c Conn, username, password, loginType, authorityType string) error {
	_, err := Send[any](ctx, c, NewLogin("global.login").
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

func GetCurrentTime(ctx context.Context, c Conn) (string, error) {
	res, err := Send[struct {
		Time string `json:"time"`
	}](ctx, c, New("global.getCurrentTime"))
	return res.Params.Time, err
}

func KeepAlive(ctx context.Context, c Conn) (int, error) {
	res, err := Send[struct {
		Timeout int `json:"timeout"`
	}](ctx, c, New("global.keepAlive"))
	return res.Params.Timeout, err
}

func Logout(ctx context.Context, c Conn) (bool, error) {
	res, err := Send[bool](ctx, c, New("global.logout"))
	return res.Params, err
}
