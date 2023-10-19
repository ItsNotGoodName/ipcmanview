package rpc

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/rpcgen"
	"github.com/rs/zerolog/log"
)

func handleErr(rpcErr rpcgen.WebRPCError, err error) error {
	log.Err(err).Msg("Request failed")
	return rpcErr
}

func convertUser(r core.User) *rpcgen.User {
	return &rpcgen.User{
		Id:        r.ID,
		Email:     r.Email,
		Username:  r.Username,
		CreatedAt: r.CreatedAt,
	}
}
