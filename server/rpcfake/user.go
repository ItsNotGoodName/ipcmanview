package rpcfake

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/server/jwt"
	"github.com/ItsNotGoodName/ipcmanview/server/rpcgen"
)

type UserService struct {
	mu        sync.Mutex
	users     []rpcgen.User
	passwords []string
	id        int
}

var _ rpcgen.AuthService = (*UserService)(nil)
var _ rpcgen.UserService = (*UserService)(nil)

func NewUserService() *UserService {
	return &UserService{}
}

// Me implements rpcgen.UserService.
func (s *UserService) Me(ctx context.Context) (*rpcgen.User, error) {
	sleep(ctx, time.Second)
	userID := jwt.DecodeUserID(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.users {
		if s.users[i].Id == userID {
			user := s.users[i]
			return &user, nil
		}
	}

	return nil, rpcgen.ErrInvalidToken
}

// Login implements rpcgen.AuthService.
func (s *UserService) Login(ctx context.Context, usernameOrEmail string, password string) (*rpcgen.User, string, error) {
	sleep(ctx, time.Second)
	usernameOrEmail = strings.ToLower(usernameOrEmail)

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.users {
		if s.users[i].Username == usernameOrEmail || s.users[i].Email == usernameOrEmail {
			if s.passwords[i] != password {
				return nil, "", rpcgen.ErrWebrpcBadRequest
			}

			user := s.users[i]

			return &user, jwt.EncodeUserID(s.users[i].Id), nil
		}
	}

	return nil, "", rpcgen.ErrWebrpcBadRequest
}

// Register implements rpcgen.AuthService.
func (s *UserService) Register(ctx context.Context, req *rpcgen.UserRegister) error {
	sleep(ctx, time.Second)
	username := strings.ToLower(req.Username)
	email := strings.ToLower(req.Email)

	if req.Password != req.PasswordConfirm {
		return rpcgen.ErrorWithCause(rpcgen.ErrWebrpcBadRequest, fmt.Errorf("password configrm does not match"))
	}
	if len(req.Password) < 3 {
		return rpcgen.ErrorWithCause(rpcgen.ErrWebrpcBadRequest, fmt.Errorf("password too short"))
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.users {
		if s.users[i].Username == username || s.users[i].Email == email {
			return rpcgen.ErrorWithCause(rpcgen.ErrWebrpcBadRequest, fmt.Errorf("user already exists"))
		}
	}

	s.id += s.id

	s.users = append(s.users, rpcgen.User{
		Id:       int64(s.id),
		Username: username,
		Email:    email,
	})
	s.passwords = append(s.passwords, req.Password)

	return nil
}
