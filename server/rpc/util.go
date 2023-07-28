package rpc

import (
	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/server/service"
)

func newUser(user core.User) *service.User {
	return &service.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}
