package dahuarpc

import (
	"context"
	"errors"
	"fmt"
)

func Login(ctx context.Context, conn ConnLogin, username, password string) error {
	firstLogin, err := FirstLogin(ctx, conn, username)
	if err != nil {
		return err
	}
	if firstLogin.Error == nil {
		return fmt.Errorf("FirstLogin did not return an error")
	}
	if !(firstLogin.Error.Code == 268632079 || firstLogin.Error.Code == 401) {
		return fmt.Errorf("FirstLogin has invalid error code: %d", firstLogin.Error.Code)
	}

	// Set session
	conn.SetSession(firstLogin.Session.String())

	// Magic
	var loginType string
	if firstLogin.Params.Encryption == "WatchNet" {
		loginType = "WatchNet"
	} else {
		loginType = "Direct"
	}

	// Encrypt password based on the first login and then do a second login
	passwordHash := firstLogin.Params.HashPassword(username, password)
	err = SecondLogin(ctx, conn, username, passwordHash, loginType, firstLogin.Params.Encryption)
	if err != nil {
		var responseErr *ResponseError
		if errors.As(err, &responseErr) {
			if loginErr := loginErrorFromResponseError(responseErr); loginErr != nil {
				return errors.Join(loginErr, err)
			}
		}

		return err
	}

	return nil
}

func loginErrorFromResponseError(r *ResponseError) *LoginError {
	switch r.Code {
	case 268632085:
		return ErrLoginUserOrPasswordNotValid
	case 268632081:
		return ErrLoginHasBeenLocked
	}

	switch r.Message {
	case "UserNotValidt":
		return ErrLoginUserNotValid
	case "PasswordNotValid":
		return ErrLoginPasswordNotValid
	case "InBlackList":
		return ErrLoginInBlackList
	case "HasBeedUsed":
		return ErrLoginHasBeedUsed
	case "HasBeenLocked":
		return ErrLoginHasBeenLocked
	}

	return nil
}

type LoginError struct {
	Message string
}

func newLoginError(message string) *LoginError {
	return &LoginError{
		Message: message,
	}
}

func (e *LoginError) Error() string {
	return e.Message
}

var (
	ErrLoginUserOrPasswordNotValid = newLoginError("User or password not valid")
	ErrLoginUserNotValid           = newLoginError("User not valid")
	ErrLoginPasswordNotValid       = newLoginError("Password not valid")
	ErrLoginInBlackList            = newLoginError("User in blackList")
	ErrLoginHasBeedUsed            = newLoginError("User has be used")
	ErrLoginHasBeenLocked          = newLoginError("User locked")
)
