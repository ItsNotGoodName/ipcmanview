package webserver

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/api"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/labstack/echo/v4"
)

type Data map[string]any

func useDahuaCamera(c echo.Context, db repo.DB) (repo.GetDahuaCameraRow, error) {
	id, err := api.PathID(c)
	if err != nil {
		return repo.GetDahuaCameraRow{}, err
	}

	camera, err := db.GetDahuaCamera(c.Request().Context(), id)
	if err != nil {
		return repo.GetDahuaCameraRow{}, err
	}

	return camera, nil
}
