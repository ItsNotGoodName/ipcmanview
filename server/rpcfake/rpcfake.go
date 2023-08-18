package rpcfake

import (
	"context"
	"math/rand"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/server/jwt"
	"github.com/ItsNotGoodName/ipcmanview/server/rpcgen"
)

func SeedAuthService(ctx context.Context, auth rpcgen.AuthService) error {
	return auth.Register(ctx, &rpcgen.UserRegister{
		Email:           "admin@example.com",
		Username:        "admin",
		Password:        "password",
		PasswordConfirm: "password",
	})
}

var _ rpcgen.AuthService = (*Service)(nil)
var _ rpcgen.UserService = (*Service)(nil)
var _ rpcgen.DahuaService = (*Service)(nil)

type Service struct {
	user rpcgen.User
}

// ActiveScanCount implements service.DahuaService.
func (*Service) ActiveScannerCount(ctx context.Context) (int, int, error) {
	return 3, 5, nil
}

// CameraCount implements service.DahuaService.
func (*Service) CameraCount(ctx context.Context) (int, error) {
	randomPanic()
	return 999, nil
}

func NewService() *Service {
	return &Service{
		user: rpcgen.User{
			Id:       3,
			Username: "Example",
			Email:    "example@example.com",
		},
	}
}

// Me implements service.UserService.
func (s *Service) Me(ctx context.Context) (*rpcgen.User, error) {
	time.Sleep(time.Second)
	return &s.user, nil
}

// Login implements service.AuthService.
func (s *Service) Login(ctx context.Context, usernameOrEmail string, password string) (*rpcgen.User, string, error) {
	// randomPanic()
	// time.Sleep(1 * time.Second)
	return &s.user, jwt.EncodeUserID(s.user.Id), nil
}

// Register implements service.AuthService.
func (*Service) Register(ctx context.Context, user *rpcgen.UserRegister) error {
	panic("unimplemented")
}

func randomPanic() {
	if rand.Int()%2 == 0 {
		panic("random panic")
	}
}

func sleep(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	case <-time.After(d):
	}
}
