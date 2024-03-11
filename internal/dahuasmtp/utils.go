package dahuasmtp

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// enableMechLogin enables the LOGIN mechanism which is used for legacy devices.
func enableMechLogin(be smtp.Backend, s *smtp.Server) {
	// Adapted from https://github.com/emersion/go-smtp/issues/41#issuecomment-493601465
	s.EnableAuth(sasl.Login, func(conn *smtp.Conn) sasl.Server {
		return sasl.NewLoginServer(func(username, password string) error {
			sess := conn.Session()
			if sess == nil {
				panic("No session when AUTH is called")
			}

			return sess.AuthPlain(username, password)
		})
	})
}
