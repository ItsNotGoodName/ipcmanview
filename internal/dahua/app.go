package dahua

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/spf13/afero"
)

var app App

type App struct {
	DB         sqlite.DB
	Hub        *bus.Hub
	AFS        afero.Fs
	Store      *Store
	ScanLocker ScanLocker
}

func Init(_app App) {
	app = _app
}
