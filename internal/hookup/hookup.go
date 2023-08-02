package hookup

import (
	"github.com/ItsNotGoodName/ipcmango/internal/dahua"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/internal/event"
	"github.com/rs/zerolog/log"
)

func DahuaStore(bus *event.Bus, store *dahua.Store) {
	bus.DahuaCameraUpdated = append(bus.DahuaCameraUpdated, func(dbCtx db.Context, evt event.DahuaCameraUpdated) {
		for _, v := range evt.IDS {
			_, err := store.GetOrCreate(dbCtx, v)
			if err != nil {
				log.Err(err).Msg("Failed to update dahua store camera")
			}
		}
	})
	bus.DahuaCameraDeleted = append(bus.DahuaCameraDeleted, func(dbCtx db.Context, evt event.DahuaCameraDeleted) {
		for _, v := range evt.IDS {
			store.Delete(dbCtx.Context, v)
		}
	})
}
