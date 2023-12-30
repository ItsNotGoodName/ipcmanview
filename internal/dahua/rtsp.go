package dahua

import "fmt"

type RTSP struct {
	Username string
	Password string
	Address  string
	Port     int
	Channel  int
	Subtype  int
}

func (r RTSP) URL() string {
	return fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/cam/realmonitor?channel=%d&subtype=%d&unicast=true&proto=Onvif",
		r.Username, r.Password, r.Address, r.Port, r.Channel, r.Subtype,
	)
}
