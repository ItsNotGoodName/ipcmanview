package mqtt

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

type Publisher struct {
	prefix   string
	address  string
	username string
	client   mqtt.Client
}

func NewPublisher(prefix, address, username, password string) Publisher {
	client := mqtt.NewClient(mqtt.
		NewClientOptions().
		AddBroker(address).
		SetUsername(username).
		SetPassword(password))
	return Publisher{
		prefix:   prefix,
		address:  address,
		username: username,
		client:   client,
	}
}

func (h Publisher) String() string {
	return "mqtt.Handler"
}

func (h Publisher) Serve(ctx context.Context) error {
	t := h.client.Connect()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.Done():
	}
	if err := t.Error(); err != nil {
		return err
	}

	log.Info().Str("address", h.address).Str("username", h.username).Msg("Connected to MQTT broker")

	<-ctx.Done()
	h.client.Disconnect(0)
	return ctx.Err()
}

func (h Publisher) Register(bus *dahua.Bus) error {
	bus.OnCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		if evt.EventRule.IgnoreMQTT {
			return nil
		}

		b, err := json.Marshal(evt.Event)
		if err != nil {
			return err
		}

		t := h.client.Publish(h.prefix+"dahua/"+strconv.FormatInt(evt.Event.CameraID, 10)+"/event", 0, false, b)
		t.Wait()
		return t.Error()
	})
	return nil
}
