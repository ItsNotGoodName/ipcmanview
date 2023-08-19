// ipcmanview-fake is used to develop the UI without the server.
package main

import (
	"context"
	"errors"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/ItsNotGoodName/ipcmanview/server"
	"github.com/ItsNotGoodName/ipcmanview/server/api"
	"github.com/ItsNotGoodName/ipcmanview/server/rpcfake"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx, shutdown := interrupt.Context()
	defer shutdown()

	svc := rpcfake.NewService()

	// Router
	r := server.Router(api.NewHandler(nil), svc, svc, svc)

	// Supervisor
	super := suture.New("root", suture.Spec{
		EventHook: sutureext.EventHook(),
	})

	// HTTP
	super.Add(server.New(r, ":8080"))

	if err := super.Serve(ctx); !errors.Is(err, context.Canceled) {
		log.Err(err).Msg("Failed to start root supervisor")
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
