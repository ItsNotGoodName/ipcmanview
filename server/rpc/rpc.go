package rpc

import (
	"github.com/ItsNotGoodName/ipcmanview/server/rpcgen"
	"github.com/rs/zerolog/log"
)

func handleErr(rpcErr rpcgen.WebRPCError, err error) error {
	log.Err(err).Msg("Request failed")
	return rpcErr
}
