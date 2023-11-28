package models

import "github.com/ItsNotGoodName/ipcmanview/internal/types"

type DTODahuaCamera struct {
	Address  string         `json:"address"`
	Location types.Location `json:"location"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Seed     int            `json:"seed"`
}
