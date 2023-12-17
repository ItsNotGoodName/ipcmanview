package mqtt

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
