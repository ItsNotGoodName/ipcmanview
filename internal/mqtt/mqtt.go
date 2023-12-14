package mqtt

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

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
		readyC:   make(chan struct{}),
	}
}

type Publisher struct {
	prefix   string
	address  string
	username string
	client   mqtt.Client
	readyC   chan struct{}
}

func (h Publisher) String() string {
	return "mqtt.Publisher"
}

func (h Publisher) Serve(ctx context.Context) error {
	select {
	case <-h.readyC:
		return suture.ErrDoNotRestart
	default:
	}

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

	close(h.readyC)
	<-ctx.Done()
	h.client.Disconnect(0)
	return ctx.Err()
}

func (h Publisher) publishDahuaEventError(ctx context.Context, cameraID int64, err error) error {
	var payload any
	if err != nil {
		payload = err.Error()
	} else {
		payload = []byte{}
	}

	t := h.client.Publish(h.prefix+"dahua/"+strconv.FormatInt(cameraID, 10)+"/event/error", 0, true, payload)
	t.Wait()
	return t.Error()
}

func (h Publisher) Register(bus *core.Bus) error {
	bus.OnDahuaCameraEvent(func(ctx context.Context, evt models.EventDahuaCameraEvent) error {
		<-h.readyC

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
	bus.OnDahuaEventWorkerConnect(func(ctx context.Context, evt models.EventDahuaEventWorkerConnect) error {
		<-h.readyC

		if err := h.publishDahuaEventError(ctx, evt.CameraID, nil); err != nil {
			return err
		}

		t := h.client.Publish(h.prefix+"dahua/"+strconv.FormatInt(evt.CameraID, 10)+"/event/state", 0, true, "online")
		t.Wait()
		return t.Error()
	})
	bus.OnDahuaEventWorkerDisconnect(func(ctx context.Context, evt models.EventDahuaEventWorkerDisconnect) error {
		<-h.readyC

		if err := h.publishDahuaEventError(ctx, evt.CameraID, evt.Error); err != nil {
			return err
		}

		{
			t := h.client.Publish(h.prefix+"dahua/"+strconv.FormatInt(evt.CameraID, 10)+"/event/state", 0, true, "offline")
			t.Wait()
			err := t.Error()
			if err != nil {
				return err
			}
		}

		return nil
	})
	return nil
}
