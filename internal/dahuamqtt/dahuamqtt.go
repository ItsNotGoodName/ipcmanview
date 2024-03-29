package dahuamqtt

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

const dahuaEventType = "dahua_event"

func newDeviceUID(deviceID string, extra ...string) string {
	if len(extra) > 0 {
		return "ipcmanview_dahua_" + deviceID + "_" + strings.Join(extra, "_")
	}
	return "ipcmanview_dahua_" + deviceID
}

func NewConn(mqtt mqtt.Conn, haEnable bool, haTopic mqtt.Topic) Conn {
	return Conn{
		conn:     mqtt,
		haEnable: haEnable,
		haTopic:  haTopic,
	}
}

type Conn struct {
	conn     mqtt.Conn
	haEnable bool
	haTopic  mqtt.Topic
}

func (Conn) String() string {
	return "dahuamqtt.Conn"
}

func (c Conn) Serve(ctx context.Context) error {
	c.conn.Ready()

	if c.haEnable {
		if err := c.haSync(ctx); err != nil {
			return err
		}
	}

	return suture.ErrDoNotRestart
}

func (c Conn) Sync(ctx context.Context) error {
	if !c.haEnable {
		return nil
	}
	c.conn.Ready()

	return c.haSync(ctx)
}

func (c Conn) haSync(ctx context.Context) error {
	c.conn.Ready()

	ids, err := dahua.ListDeviceIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		if err := c.haSyncDevice(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (c Conn) haSyncDevice(ctx context.Context, id int64) error {
	device, err := dahua.GetDevice(ctx, id)
	if err != nil {
		return err
	}

	client, err := dahua.GetClient(ctx, id)
	if err != nil {
		return err
	}

	detail, err := dahua.GetDahuaDetail(ctx, client.RPC)
	if err != nil {
		log.Err(err).Msg("Failed to get detail")
		return nil
	}

	sw, err := dahua.GetSoftwareVersion(ctx, client.RPC)
	if err != nil {
		log.Err(err).Msg("Failed to get software version")
		return nil
	}

	coaxialCaps, err := dahua.GetCoaxialCaps(ctx, client.RPC, 1)
	if err != nil {
		log.Err(err).Msg("Failed to get coaxial caps")
		return nil
	}

	deviceID := mqtt.Int(device.ID)
	deviceUID := newDeviceUID(deviceID)

	haEntity := mqtt.NewHaEntity(c.conn)
	haEntity.Device.Name = device.Name
	haEntity.Device.Manufacturer = detail.Vendor
	haEntity.Device.Model = detail.DeviceType
	haEntity.Device.HwVersion = detail.HardwareVersion
	haEntity.Device.SwVersion = sw.Version
	haEntity.Device.Identifiers = []string{deviceUID}
	haEntity.ObjectId = "dahua_" + device.Name

	// event
	{
		topicDahuaIDEvent := mqtt.Topic(c.conn.Topic.Join("dahua", deviceID, "event"))

		event := mqtt.HaEvent{HaEntity: haEntity}
		event.Availability = append(event.Availability, mqtt.HaAvailability{
			Topic: topicDahuaIDEvent.Join("state"),
		})
		event.StateTopic = string(topicDahuaIDEvent)
		event.UniqueId = deviceUID
		event.Name = "Event"
		event.EventTypes = []string{dahuaEventType}

		b, err := json.Marshal(event)
		if err != nil {
			return err
		}

		topicConfig := c.haTopic.Join("event", deviceUID, "config")
		if err := mqtt.Wait(c.conn.Client.Publish(topicConfig, 0, true, b)); err != nil {
			return err
		}
	}

	// white_light
	if coaxialCaps.SupportControlLight {
		topicDahuaIDWhiteLight := mqtt.Topic(c.conn.Topic.Join("dahua", deviceID, "white_light"))

		binarySensor := mqtt.HaBinarySensor{HaEntity: haEntity}
		binarySensor.StateTopic = string(topicDahuaIDWhiteLight)
		binarySensor.UniqueId = newDeviceUID(deviceID, "white_light")
		binarySensor.Name = "White Light"
		binarySensor.Icon = "mdi:lightbulb"

		b, err := json.Marshal(binarySensor)
		if err != nil {
			return err
		}

		topicConfig := c.haTopic.Join("binary_sensor", deviceUID, "white_light", "config")
		if err := mqtt.Wait(c.conn.Client.Publish(topicConfig, 0, true, b)); err != nil {
			return err
		}
	}

	// speaker
	if coaxialCaps.SupportControlSpeaker {
		topicDahuaIDSpeaker := mqtt.Topic(c.conn.Topic.Join("dahua", deviceID, "speaker"))

		binarySensor := mqtt.HaBinarySensor{HaEntity: haEntity}
		binarySensor.StateTopic = string(topicDahuaIDSpeaker)
		binarySensor.UniqueId = newDeviceUID(deviceID, "speaker")
		binarySensor.Name = "Speaker"
		binarySensor.Icon = "mdi:bullhorn"

		b, err := json.Marshal(binarySensor)
		if err != nil {
			return err
		}

		topicConfig := c.haTopic.Join("binary_sensor", deviceUID, "speaker", "config")
		if err := mqtt.Wait(c.conn.Client.Publish(topicConfig, 0, true, b)); err != nil {
			return err
		}
	}

	return nil
}

type Event struct {
	ID        int64           `json:"id"`
	DeviceID  int64           `json:"device_id"`
	Code      string          `json:"code"`
	Action    string          `json:"action"`
	Index     int             `json:"index"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	EventType string          `json:"event_type"`
}

func NewEvent(v repo.DahuaEvent) Event {
	return Event{
		ID:        v.ID,
		DeviceID:  v.DeviceID,
		Code:      v.Code,
		Action:    v.Action,
		Index:     int(v.Index),
		Data:      v.Data.RawMessage,
		CreatedAt: v.CreatedAt.Time,
		EventType: dahuaEventType,
	}
}

func (c Conn) Register(hub *bus.Hub) Conn {
	if c.haEnable {
		hub.OnDahuaDeviceCreated(c.String(), func(ctx context.Context, event bus.DahuaDeviceCreated) error {
			c.conn.Ready()
			return c.haSyncDevice(ctx, event.DeviceID)
		})
		hub.OnDahuaDeviceUpdated(c.String(), func(ctx context.Context, event bus.DahuaDeviceUpdated) error {
			c.conn.Ready()
			return c.haSyncDevice(ctx, event.DeviceID)
		})
	}
	hub.OnDahuaEvent(c.String(), func(ctx context.Context, event bus.DahuaEvent) error {
		c.conn.Ready()

		if event.EventRule.IgnoreMqtt {
			return nil
		}

		b, err := json.Marshal(NewEvent(event.Event))
		if err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.Event.DeviceID), "event"), 0, false, b))
	})
	hub.OnDahuaWorkerConnected(c.String(), func(ctx context.Context, event bus.DahuaWorkerConnected) error {
		c.conn.Ready()

		if err := publishDeviceError(ctx, c.conn, event.DeviceID, string(event.Type), nil); err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", strconv.FormatInt(event.DeviceID, 10), string(event.Type), "state"), 0, true, "online"))
	})
	hub.OnDahuaWorkerDisconnected(c.String(), func(ctx context.Context, event bus.DahuaWorkerDisconnected) error {
		c.conn.Ready()

		if err := publishDeviceError(ctx, c.conn, event.DeviceID, string(event.Type), event.Error); err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.DeviceID), string(event.Type), "state"), 0, true, "offline"))
	})
	hub.OnDahuaCoaxialStatus(c.String(), func(ctx context.Context, event bus.DahuaCoaxialStatus) error {
		c.conn.Ready()

		{
			payload := "OFF"
			if event.CoaxialStatus.WhiteLight {
				payload = "ON"
			}

			if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.DeviceID), string(models.DahuaWorkerType_Coaxial), "white_light"), 0, true, payload)); err != nil {
				return err
			}
		}

		{
			payload := "OFF"
			if event.CoaxialStatus.Speaker {
				payload = "ON"
			}

			if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.DeviceID), string(models.DahuaWorkerType_Coaxial), "speaker"), 0, true, payload)); err != nil {
				return err
			}
		}

		return nil
	})
	hub.OnDahuaFileCursorUpdated(c.String(), func(ctx context.Context, event bus.DahuaFileCursorUpdated) error {
		c.conn.Ready()

		payload, err := json.Marshal(event.Cursor)
		if err != nil {
			return err
		}

		if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.Cursor.DeviceID), "cursor"), 0, true, payload)); err != nil {
			return err
		}

		return nil
	})
	return c
}

func publishDeviceError(ctx context.Context, conn mqtt.Conn, deviceID int64, resource string, err error) error {
	var payload any
	if err != nil {
		payload = err.Error()
	} else {
		payload = []byte{}
	}
	return mqtt.Wait(conn.Client.Publish(conn.Topic.Join("dahua", mqtt.Int(deviceID), resource, "error"), 0, true, payload))
}
