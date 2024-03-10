package dahuatasks

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

var app App

type App struct {
	DB  sqlite.DB
	Hub *bus.Hub
}

func Init(_app App) {
	app = _app
}
