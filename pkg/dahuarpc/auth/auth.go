package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/global"
)

const (
	WatchNet = "WatchNet"
	TimeOut  = 60 * time.Second
)

func Logout(ctx context.Context, conn *dahuarpc.Conn) {
	global.Logout(ctx, conn)
	conn.Set(dahuarpc.StateLogout)
}

func KeepAlive(ctx context.Context, conn *dahuarpc.Conn) error {
	if time.Now().Sub(conn.LastLogin) > TimeOut {
		_, err := global.KeepAlive(ctx, conn)
		if err != nil {
			if !errors.Is(err, dahuarpc.ErrRequestFailed) {
				conn.Set(dahuarpc.StateLogout)
			}

			return err
		}

		conn.Set(dahuarpc.StateLogin)
	}

	return nil
}

func Login(ctx context.Context, conn *dahuarpc.Conn, username, password string) error {
	if err := login(ctx, conn, username, password); err != nil {
		var e *LoginError
		if errors.As(err, &e) {
			conn.Set(dahuarpc.StateError, err)
		} else {
			conn.Set(dahuarpc.StateLogout)
		}

		return err
	}

	conn.Set(dahuarpc.StateLogin)

	return nil
}

func login(ctx context.Context, conn *dahuarpc.Conn, username, password string) error {
	// Do a first login
	firstLogin, err := global.FirstLogin(ctx, conn, username)
	if err != nil {
		return err
	}
	if firstLogin.Error == nil {
		return fmt.Errorf("FirstLogin did not return an error")
	}
	if !(firstLogin.Error.Code == 268632079 || firstLogin.Error.Code == 401) {
		return fmt.Errorf("FirstLogin has invalid error code: %d", firstLogin.Error.Code)
	}

	// Update session
	conn.UpdateSession(firstLogin.Session.String())

	// Magic
	loginType := func() string {
		if firstLogin.Params.Encryption == WatchNet {
			return WatchNet
		}
		return "Direct"
	}()

	// Encrypt password based on the first login and then do a second login
	passwordHash := firstLogin.Params.HashPassword(username, password)
	err = global.SecondLogin(ctx, conn, username, passwordHash, loginType, firstLogin.Params.Encryption)
	if err != nil {
		var responseErr *dahuarpc.ErrResponse
		if errors.As(err, &responseErr) {
			if loginErr := intoLoginError(responseErr); loginErr != nil {
				return errors.Join(loginErr, err)
			}
		}

		return err
	}

	return nil
}

func intoLoginError(r *dahuarpc.ErrResponse) *LoginError {
	switch r.Code {
	case 268632085:
		return &ErrLoginUserOrPasswordNotValid
	case 268632081:
		return &ErrLoginHasBeenLocked
	}

	switch r.Message {
	case "UserNotValidt":
		return &ErrLoginUserNotValid
	case "PasswordNotValid":
		return &ErrLoginPasswordNotValid
	case "InBlackList":
		return &ErrLoginInBlackList
	case "HasBeedUsed":
		return &ErrLoginHasBeedUsed
	case "HasBeenLocked":
		return &ErrLoginHasBeenLocked
	}

	return nil
}

type LoginError struct {
	Message string
}

func newErrLogin(message string) LoginError {
	return LoginError{
		Message: message,
	}
}

func (e *LoginError) Error() string {
	return e.Message
}

var (
	ErrLoginClosed                 = newErrLogin("Client is closed")
	ErrLoginUserOrPasswordNotValid = newErrLogin("User or password not valid")
	ErrLoginUserNotValid           = newErrLogin("User not valid")
	ErrLoginPasswordNotValid       = newErrLogin("Password not valid")
	ErrLoginInBlackList            = newErrLogin("User in blackList")
	ErrLoginHasBeedUsed            = newErrLogin("User has be used")
	ErrLoginHasBeenLocked          = newErrLogin("User locked")
)
