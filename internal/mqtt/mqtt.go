package mqtt

import (
	"context"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func Int(i int64) string {
	return strconv.FormatInt(i, 10)
}

type Topic string

func (t Topic) Join(topics ...string) string {
	return string(t) + "/" + Join(topics...)
}

func Join(topics ...string) string {
	return strings.Join(topics, "/")
}

func Wait(t mqtt.Token) error {
	t.Wait()
	return t.Error()
}

func TopicServerState(topic Topic) string {
	return topic.Join("server", "state")
}

func NewConn(topic Topic, address, username, password string) Conn {
	client := mqtt.NewClient(mqtt.NewClientOptions().
		AddBroker(address).
		SetUsername(username).
		SetPassword(password).
		SetWill(TopicServerState(topic), "offline", 0, true).
		SetOnConnectHandler(func(c mqtt.Client) {
			c.Publish(TopicServerState(topic), 0, true, "online")
		}))
	return Conn{
		Client:   client,
		Topic:    topic,
		address:  address,
		username: username,
		readyC:   make(chan struct{}),
	}
}

type Conn struct {
	Client   mqtt.Client
	Topic    Topic
	address  string
	username string
	readyC   chan struct{}
}

func (h Conn) String() string {
	return "mqtt.Conn"
}

func (h Conn) Serve(ctx context.Context) error {
	select {
	case <-h.readyC:
		return suture.ErrDoNotRestart
	default:
	}

	t := h.Client.Connect()
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
	h.Client.Disconnect(0)
	return ctx.Err()
}

func (h Conn) Ready() {
	<-h.readyC
}
