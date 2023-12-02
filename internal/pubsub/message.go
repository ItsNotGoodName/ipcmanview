package pubsub

import (
	"encoding/json"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/mochi-mqtt/server/v2/packets"
)

func SendDahuaEvent(server *Pub, event models.DahuaEvent) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return server.Publish("dahua/"+strconv.FormatInt(event.CameraID, 10)+"/event", b, false, 0)
}

func ParseDahuaEvent(pk packets.Packet) (models.DahuaEvent, error) {
	var event models.DahuaEvent
	err := json.Unmarshal(pk.Payload, &event)
	return event, err
}
