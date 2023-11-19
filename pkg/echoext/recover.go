package echoext

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func RecoverLogErrorFunc(c echo.Context, err error, stack []byte) error {
	log.Err(err).Msgf("%s", stack)
	return nil
}
