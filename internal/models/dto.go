package models

type DTODahuaCamera struct {
	Address  string   `json:"address"`
	Location Location `json:"location"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Seed     int      `json:"seed"`
}
