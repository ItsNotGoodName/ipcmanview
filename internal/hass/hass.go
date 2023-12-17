package hass

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/thejerf/suture/v4"
)

type Conn struct {
	conn  mqtt.Conn
	Topic mqtt.Topic
	db    repo.DB
}

func NewConn(mqtt mqtt.Conn, db repo.DB, topic mqtt.Topic) Conn {
	return Conn{
		conn:  mqtt,
		Topic: topic,
		db:    db,
	}
}

func (ha Conn) Serve(ctx context.Context) error {
	ha.conn.Ready()

	cameras, err := ha.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	for _, camera := range cameras {
		cameraEventTopic := mqtt.Topic(ha.conn.Topic.Join("dahua", strconv.FormatInt(camera.ID, 10), "event"))
		cfg := NewConfig(ha.conn)
		cfg.Device.Name = camera.Name
		cfg.Device.Manufacturer = "Dahua"
		cfg.Device.Identifiers = []string{strconv.FormatInt(camera.ID, 10)}
		cfg.Availability = append(cfg.Availability, ConfigAvailability{
			Topic: cameraEventTopic.Join("state"),
		})
		cfg.StateTopic = string(cameraEventTopic)
		cfg.EventTypes = []string{mqtt.DahuaEventType}
		cfg.UniqueID = "ipcmanview_dahua_" + strconv.FormatInt(camera.ID, 10)

		b, err := json.Marshal(cfg)
		if err != nil {
			return err
		}

		topic := ha.Topic.Join("event", "ipcmanview-dahua-"+strconv.FormatInt(camera.ID, 10), "config")
		if err := mqtt.Wait(ha.conn.Client.Publish(topic, 0, true, b)); err != nil {
			return err
		}
	}

	return suture.ErrDoNotRestart
}
