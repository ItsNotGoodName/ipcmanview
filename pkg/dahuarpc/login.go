package dahuarpc

import (
	"context"
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
		return err
	}

	return nil
}
