package core

import (
	"crypto/rand"
	"encoding/base64"
)

// RuntimeToken is a token generated at runtime to authenticate the application with itself.
var RuntimeToken string

func init() {
	b := make([]byte, 64)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	RuntimeToken = base64.URLEncoding.EncodeToString(b)
}
