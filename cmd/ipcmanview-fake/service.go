package main

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/server/jwt"
	"github.com/ItsNotGoodName/ipcmanview/server/service"
)

var _ service.AuthService = (*Service)(nil)
var _ service.UserService = (*Service)(nil)
var _ service.DahuaService = (*Service)(nil)

type Service struct {
	user service.User
}

// CameraCount implements service.DahuaService.
func (*Service) CameraCount(ctx context.Context) (int, error) {
	randomPanic()
	return 999, nil
}

func NewService() *Service {
	return &Service{
		user: service.User{
			Id:       1,
			Username: "Example",
			Email:    "example@example.com",
		},
	}
}

// Me implements service.UserService.
func (s *Service) Me(ctx context.Context) (*service.User, error) {
	return &s.user, nil
}

// Login implements service.AuthService.
func (s *Service) Login(ctx context.Context, usernameOrEmail string, password string) (*service.User, string, error) {
	return &s.user, jwt.EncodeUserID(s.user.Id), nil
}

// Register implements service.AuthService.
func (*Service) Register(ctx context.Context, user *service.UserRegister) error {
	panic("unimplemented")
}
