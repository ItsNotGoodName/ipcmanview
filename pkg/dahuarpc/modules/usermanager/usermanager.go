package usermanager

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type User struct {
	ClientAddress string             `json:"ClientAddress"`
	ClientType    string             `json:"ClientType"`
	Group         string             `json:"Group"`
	ID            int                `json:"Id"`
	LoginTime     dahuarpc.Timestamp `json:"LoginTime"`
	Name          string             `json:"Name"`
}

func GetActiveUserInfoAll(ctx context.Context, c dahuarpc.Conn) ([]User, error) {
	res, err := dahuarpc.Send[struct {
		Users []User `json:"users"`
	}](ctx, c, dahuarpc.
		New("userManager.getActiveUserInfoAll"))
	if err != nil {
		return nil, err
	}

	return res.Params.Users, nil
}

func GetAuthorityList(ctx context.Context, c dahuarpc.Conn) ([]string, error) {
	res, err := dahuarpc.Send[[]string](ctx, c, dahuarpc.
		New("userManager.getAuthorityList"))
	if err != nil {
		return nil, err
	}

	return res.Params, nil
}

type UserInfo struct {
	Anonymous            bool     `json:"Anonymous"`
	AuthorityList        []string `json:"AuthorityList"`
	Group                string   `json:"Group"`
	ID                   int      `json:"Id"`
	Memo                 string   `json:"Memo"`
	Name                 string   `json:"Name"`
	Password             string   `json:"Password"`
	PasswordModifiedTime string   `json:"PasswordModifiedTime"`
	PwdScore             int      `json:"PwdScore"`
	Reserved             bool     `json:"Reserved"`
	Sharable             bool     `json:"Sharable"`
}

func GetUserInfoAll(ctx context.Context, c dahuarpc.Conn) ([]UserInfo, error) {
	res, err := dahuarpc.Send[struct {
		Users []UserInfo `json:"users"`
	}](ctx, c, dahuarpc.
		New("userManager.getUserInfoAll"))
	if err != nil {
		return nil, err
	}

	return res.Params.Users, nil
}

type GroupInfo struct {
	AuthorityList []string `json:"AuthorityList"`
	ID            int      `json:"Id"`
	Memo          string   `json:"Memo"`
	Name          string   `json:"Name"`
}

func GetGroupInfoAll(ctx context.Context, c dahuarpc.Conn) ([]GroupInfo, error) {
	res, err := dahuarpc.Send[[]GroupInfo](ctx, c, dahuarpc.
		New("userManager.getGroupInfoAll"))
	if err != nil {
		return nil, err
	}

	return res.Params, nil
}
