package auth

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
)

var app App

type App struct {
	DB                   sqlite.DB
	Hub                  *bus.Hub
	TouchSessionThrottle TouchSessionThrottle
}

func Init(_app App) {
	app = _app
}
