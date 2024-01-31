package dahuamqtt

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
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

type Conn struct {
	conn     mqtt.Conn
	db       repo.DB
	store    *dahua.Store
	haEnable bool
	haTopic  mqtt.Topic
}

func NewConn(mqtt mqtt.Conn, db repo.DB, store *dahua.Store, haEnable bool, haTopic mqtt.Topic) Conn {
	return Conn{
		conn:     mqtt,
		db:       db,
		store:    store,
		haEnable: haEnable,
		haTopic:  haTopic,
	}
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

	devices, err := c.db.DahuaListFatDevices(ctx)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if err := c.haSyncDevice(ctx, device.DahuaDevice, dahua.NewConn(device)); err != nil {
			return err
		}
	}

	return nil
}

func (c Conn) haSyncDevice(ctx context.Context, device repo.DahuaDevice, conn dahua.Conn) error {
	client := c.store.Client(ctx, conn)

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
		Data:      v.Data,
		CreatedAt: v.CreatedAt.Time,
		EventType: dahuaEventType,
	}
}

func (c Conn) Register(bus *event.Bus) error {
	if c.haEnable {
		bus.OnDahuaDeviceCreated(func(ctx context.Context, evt event.DahuaDeviceCreated) error {
			c.conn.Ready()

			return c.haSyncDevice(ctx, evt.Device.DahuaDevice, dahua.NewConn(evt.Device))
		})
		bus.OnDahuaDeviceUpdated(func(ctx context.Context, evt event.DahuaDeviceUpdated) error {
			c.conn.Ready()

			return c.haSyncDevice(ctx, evt.Device.DahuaDevice, dahua.NewConn(evt.Device))
		})
	}
	bus.OnDahuaEvent(func(ctx context.Context, evt event.DahuaEvent) error {
		c.conn.Ready()

		if evt.EventRule.IgnoreMqtt {
			return nil
		}

		b, err := json.Marshal(NewEvent(evt.Event))
		if err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(evt.Event.DeviceID), "event"), 0, false, b))
	})
	bus.OnDahuaEventWorkerConnect(func(ctx context.Context, evt event.DahuaEventWorkerConnect) error {
		c.conn.Ready()

		if err := publishEventError(ctx, c.conn, evt.DeviceID, nil); err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", strconv.FormatInt(evt.DeviceID, 10), "event", "state"), 0, true, "online"))
	})
	bus.OnDahuaEventWorkerDisconnect(func(ctx context.Context, evt event.DahuaEventWorkerDisconnect) error {
		c.conn.Ready()

		if err := publishEventError(ctx, c.conn, evt.DeviceID, evt.Error); err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(evt.DeviceID), "event", "state"), 0, true, "offline"))
	})
	bus.OnDahuaCoaxialStatus(func(ctx context.Context, event event.DahuaCoaxialStatus) error {
		c.conn.Ready()

		{
			payload := "OFF"
			if event.CoaxialStatus.WhiteLight {
				payload = "ON"
			}

			if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.DeviceID), "white_light"), 0, true, payload)); err != nil {
				return err
			}
		}

		{
			payload := "OFF"
			if event.CoaxialStatus.Speaker {
				payload = "ON"
			}

			if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.DeviceID), "speaker"), 0, true, payload)); err != nil {
				return err
			}
		}

		return nil
	})
	return nil
}

func publishEventError(ctx context.Context, conn mqtt.Conn, deviceID int64, err error) error {
	var payload any
	if err != nil {
		payload = err.Error()
	} else {
		payload = []byte{}
	}
	return mqtt.Wait(conn.Client.Publish(conn.Topic.Join("dahua", mqtt.Int(deviceID), "event", "error"), 0, true, payload))
}
