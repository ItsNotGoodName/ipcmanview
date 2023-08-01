package procs

import (
	"github.com/ItsNotGoodName/ipcmango/internal/dahua"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/internal/event"
)

func HookBusToDahuaStore(bus *event.Bus, store *dahua.Store) {
	bus.DahuaCameraDeleted = append(bus.DahuaCameraDeleted, func(context db.Context, evt event.DahuaCameraDeleted) {
		for _, v := range evt.IDS {
			store.Delete(context.Context, v)
		}
	})
}
