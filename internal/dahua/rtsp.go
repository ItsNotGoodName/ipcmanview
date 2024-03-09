package dahua

import "fmt"

type GetLiveRTSPURLParams struct {
	Username string
	Password string
	Host     string
	Port     int
	Channel  int
	Subtype  int
}

func GetLiveRTSPURL(arg GetLiveRTSPURLParams) string {
	return fmt.Sprintf(
		"rtsp://%s:%s@%s:%d/cam/realmonitor?channel=%d&subtype=%d&unicast=true&proto=Onvif",
		arg.Username, arg.Password, arg.Host, arg.Port, arg.Channel, arg.Subtype,
	)
}
