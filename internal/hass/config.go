package hass

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
)

// homeassistant/event/dahua-camera-1/config
// {
//   "availability": [
//     {
//       "topic": "ipcmanview/server/state"
//     },
//     {
//       "topic": "ipcmanview/dahua/21/event/state"
//     }
//   ],
//   "availability_mode": "all",
//   "event_types": [
//     "dahua_event"
//   ],
//   "state_topic": "ipcmanview/dahua/21/event"
// }

func NewConfig(conn mqtt.Conn) Config {
	return Config{
		Availability: []ConfigAvailability{
			{Topic: mqtt.TopicServerState(conn.Topic)},
		},
		AvailabilityMode: "all",
		Origin: ConfigOrigin{
			Name:      "IPCManView",
			SWVersion: build.Current.Version,
			URL:       build.Current.RepoURL,
		},
	}
}

type Config struct {
	Availability     []ConfigAvailability `json:"availability,omitempty"`
	AvailabilityMode string               `json:"availability_mode,omitempty"`
	Device           ConfigDevice         `json:"device,omitempty"`
	EventTypes       []string             `json:"event_types,omitempty"`
	Origin           ConfigOrigin         `json:"origin,omitempty"`
	Name             string               `json:"name,omitempty"`
	ObjectID         string               `json:"object_id,omitempty"`
	StateTopic       string               `json:"state_topic,omitempty"`
	UniqueID         string               `json:"unique_id,omitempty"`
}

type ConfigAvailability struct {
	Topic string `json:"topic,omitempty"`
}

type ConfigDevice struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Name         string   `json:"name,omitempty"`
	SWVersion    string   `json:"sw_version,omitempty"`
}

type ConfigOrigin struct {
	Name      string `json:"name,omitempty"`
	SWVersion string `json:"sw_version,omitempty"`
	URL       string `json:"url,omitempty"`
}
