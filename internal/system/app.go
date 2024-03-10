package system

var app App

type App struct {
	CP ConfigProvider
}

func Init(_app App) {
	app = _app
}
