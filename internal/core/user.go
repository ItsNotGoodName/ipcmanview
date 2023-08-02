package core

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}

func (u User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type UserCreate struct {
	Username        string `validate:"gte=3,lte=64"`
	Email           string `validate:"required,email"`
	Password        string `validate:"gte=8"`
	PasswordConfirm string
}

func NewUser(r UserCreate) (User, error) {
	if r.Password != r.PasswordConfirm {
		return User{}, fmt.Errorf("passwords do not match")
	}

	if err := validate.Struct(&r); err != nil {
		return User{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	email := strings.ToLower(r.Email)

	return User{
		ID:       0,
		Email:    email,
		Username: r.Username,
		Password: string(password),
	}, nil
}
