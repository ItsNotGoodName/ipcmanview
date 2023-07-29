package client

import (
	"fmt"

	"github.com/ItsNotGoodName/pkg/dahua"
	"github.com/ItsNotGoodName/pkg/dahua/modules/global"
)

const Timeout = 60
const WatchNet = "WatchNet"

func Logout(conn *dahua.Conn, gen dahua.Generator) {
	global.Logout(gen)
	conn.Set(dahua.StateLogout)
}

func Login(conn *dahua.Conn, gen dahua.Generator, username, password string) error {
	if conn.State == dahua.StateLogin {
		Logout(conn, gen)
	} else if conn.State == dahua.StateError {
		panic("error")
	}

	err := login(conn, gen, username, password)
	if err != nil {
		panic("TODO: set state to error")
	}

	conn.Set(dahua.StateLogin)

	return nil
}

func login(conn *dahua.Conn, gen dahua.Generator, username, password string) error {
	// Do a first login
	firstLogin, err := global.FirstLogin(gen, username)
	if err != nil {
		return err
	}
	if firstLogin.Error == nil {
		panic("Error was not supposed to be nil")
	}
	if !(firstLogin.Error.Code == 268632079 || firstLogin.Error.Code == 401) {
		panic(fmt.Sprintf("invalid error code %d", firstLogin.Error.Code))
	}

	// Set session
	if err := conn.SetSession(firstLogin.Session.Value); err != nil {
		return err
	}

	// Magic
	loginType := func() string {
		if firstLogin.Params.Encryption == WatchNet {
			return WatchNet
		}
		return "Direct"
	}()

	// Encrypt password based on the first login and then do a second login
	passwordHash := firstLogin.Params.HashPassword(username, password)
	err = global.SecondLogin(gen, username, passwordHash, loginType, firstLogin.Params.Encryption)
	if err != nil {
		panic(fmt.Sprintf("I can't handle this: %s", err))
	}

	return nil
}
