package dahuamqtt

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/mqtt"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

const dahuaEventType = "dahua_event"

func newCameraUID(cameraID string, extra ...string) string {
	if len(extra) > 0 {
		return "ipcmanview_dahua_" + cameraID + "_" + strings.Join(extra, "_")
	}
	return "ipcmanview_dahua_" + cameraID
}

type Conn struct {
	conn     mqtt.Conn
	db       repo.DB
	store    *dahuacore.Store
	haEnable bool
	haTopic  mqtt.Topic
}

func NewConn(mqtt mqtt.Conn, db repo.DB, store *dahuacore.Store, haEnable bool, haTopic mqtt.Topic) Conn {
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

	cameras, err := c.db.ListDahuaCamera(ctx)
	if err != nil {
		return err
	}

	for _, dbCamera := range cameras {
		if err := c.haSyncCamera(ctx, dbCamera.Convert()); err != nil {
			return err
		}
	}

	return nil
}

func (c Conn) haSyncCamera(ctx context.Context, camera models.DahuaCameraConn) error {
	conn := c.store.Conn(ctx, camera.DahuaConn)

	detail, err := dahuacore.GetDahuaDetail(ctx, conn.Camera.ID, conn.RPC)
	if err != nil {
		log.Err(err).Msg("Failed to get detail")
		return nil
	}

	sw, err := dahuacore.GetSoftwareVersion(ctx, conn.Camera.ID, conn.RPC)
	if err != nil {
		log.Err(err).Msg("Failed to get software version")
		return nil
	}

	coaxialCaps, err := dahuacore.GetCoaxialCaps(ctx, conn.Camera.ID, conn.RPC, 1)
	if err != nil {
		log.Err(err).Msg("Failed to get coaxial caps")
		return nil
	}

	cameraID := mqtt.Int(camera.DahuaCamera.ID)
	cameraUID := newCameraUID(cameraID)

	haEntity := mqtt.NewHaEntity(c.conn)
	haEntity.Device.Name = camera.Name
	haEntity.Device.Manufacturer = detail.Vendor
	haEntity.Device.Model = detail.DeviceType
	haEntity.Device.HwVersion = detail.HardwareVersion
	haEntity.Device.SwVersion = sw.Version
	haEntity.Device.Identifiers = []string{cameraUID}
	haEntity.ObjectId = "dahua_" + camera.Name

	// event
	{
		topicDahuaIDEvent := mqtt.Topic(c.conn.Topic.Join("dahua", cameraID, "event"))

		event := mqtt.HaEvent{HaEntity: haEntity}
		event.Availability = append(event.Availability, mqtt.HaAvailability{
			Topic: topicDahuaIDEvent.Join("state"),
		})
		event.StateTopic = string(topicDahuaIDEvent)
		event.UniqueId = cameraUID
		event.Name = "Event"
		event.EventTypes = []string{dahuaEventType}

		b, err := json.Marshal(event)
		if err != nil {
			return err
		}

		topicConfig := c.haTopic.Join("event", cameraUID, "config")
		if err := mqtt.Wait(c.conn.Client.Publish(topicConfig, 0, true, b)); err != nil {
			return err
		}
	}

	// white_light
	if coaxialCaps.SupportControlLight {
		topicDahuaIDWhiteLight := mqtt.Topic(c.conn.Topic.Join("dahua", cameraID, "white_light"))

		binarySensor := mqtt.HaBinarySensor{HaEntity: haEntity}
		binarySensor.StateTopic = string(topicDahuaIDWhiteLight)
		binarySensor.UniqueId = newCameraUID(cameraID, "white_light")
		binarySensor.Name = "White Light"
		binarySensor.Icon = "mdi:lightbulb"

		b, err := json.Marshal(binarySensor)
		if err != nil {
			return err
		}

		topicConfig := c.haTopic.Join("binary_sensor", cameraUID, "white_light", "config")
		if err := mqtt.Wait(c.conn.Client.Publish(topicConfig, 0, true, b)); err != nil {
			return err
		}
	}

	// speaker
	if coaxialCaps.SupportControlSpeaker {
		topicDahuaIDSpeaker := mqtt.Topic(c.conn.Topic.Join("dahua", cameraID, "speaker"))

		binarySensor := mqtt.HaBinarySensor{HaEntity: haEntity}
		binarySensor.StateTopic = string(topicDahuaIDSpeaker)
		binarySensor.UniqueId = newCameraUID(cameraID, "speaker")
		binarySensor.Name = "Speaker"
		binarySensor.Icon = "mdi:bullhorn"

		b, err := json.Marshal(binarySensor)
		if err != nil {
			return err
		}

		topicConfig := c.haTopic.Join("binary_sensor", cameraUID, "speaker", "config")
		if err := mqtt.Wait(c.conn.Client.Publish(topicConfig, 0, true, b)); err != nil {
			return err
		}
	}

	return nil
}

type dahuaEvent struct {
	models.DahuaEvent
	EventType string `json:"event_type"`
}

func (c Conn) Register(bus *core.Bus) error {
	if c.haEnable {
		bus.OnEventDahuaCameraCreated(func(ctx context.Context, event models.EventDahuaCameraCreated) error {
			c.conn.Ready()

			return c.haSyncCamera(ctx, event.Camera)
		})
		bus.OnEventDahuaCameraUpdated(func(ctx context.Context, event models.EventDahuaCameraUpdated) error {
			c.conn.Ready()

			return c.haSyncCamera(ctx, event.Camera)
		})
	}
	bus.OnEventDahuaCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		c.conn.Ready()

		if evt.EventRule.IgnoreMQTT {
			return nil
		}

		b, err := json.Marshal(dahuaEvent{DahuaEvent: evt.Event, EventType: dahuaEventType})
		if err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(evt.Event.CameraID), "event"), 0, false, b))
	})
	bus.OnEventDahuaEventWorkerConnect(func(ctx context.Context, evt models.EventDahuaEventWorkerConnect) error {
		c.conn.Ready()

		if err := publishEventError(ctx, c.conn, evt.CameraID, nil); err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", strconv.FormatInt(evt.CameraID, 10), "event", "state"), 0, true, "online"))
	})
	bus.OnEventDahuaEventWorkerDisconnect(func(ctx context.Context, evt models.EventDahuaEventWorkerDisconnect) error {
		c.conn.Ready()

		if err := publishEventError(ctx, c.conn, evt.CameraID, evt.Error); err != nil {
			return err
		}

		return mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(evt.CameraID), "event", "state"), 0, true, "offline"))
	})
	bus.OnEventDahuaCoaxialStatus(func(ctx context.Context, event models.EventDahuaCoaxialStatus) error {
		c.conn.Ready()

		{
			payload := "OFF"
			if event.CoaxialStatus.WhiteLight {
				payload = "ON"
			}

			if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.CoaxialStatus.CameraID), "white_light"), 0, true, payload)); err != nil {
				return err
			}
		}

		{
			payload := "OFF"
			if event.CoaxialStatus.Speaker {
				payload = "ON"
			}

			if err := mqtt.Wait(c.conn.Client.Publish(c.conn.Topic.Join("dahua", mqtt.Int(event.CoaxialStatus.CameraID), "speaker"), 0, true, payload)); err != nil {
				return err
			}
		}

		return nil
	})
	return nil
}

func publishEventError(ctx context.Context, conn mqtt.Conn, cameraID int64, err error) error {
	var payload any
	if err != nil {
		payload = err.Error()
	} else {
		payload = []byte{}
	}
	return mqtt.Wait(conn.Client.Publish(conn.Topic.Join("dahua", mqtt.Int(cameraID), "event", "error"), 0, true, payload))
}
