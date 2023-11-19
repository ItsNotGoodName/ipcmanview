package models

type DTODahuaCamera struct {
	Address  string   `json:"address"`
	Location Location `json:"location"`
	Password string   `json:"password"`
	Username string   `json:"username"`
}
