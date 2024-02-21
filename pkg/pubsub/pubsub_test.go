package pubsub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func noopEventHandler(ctx context.Context, event Event) error { return nil }

func TestPub(t *testing.T) {
	ctx := context.Background()

	pub := NewPub()

	topics := []string{"Potato", "Thing"}
	eventTopics := func() []Event {
		eventTopics := make([]Event, 0, len(topics))
		for _, t := range topics {
			eventTopics = append(eventTopics, EventTopic(t))
		}
		return eventTopics
	}()

	sub, err := pub.
		Subscribe(eventTopics...).
		Function(ctx, noopEventHandler)
	if !assert.NoError(t, err) {
		return
	}
	defer sub.Close()

	s, err := pub.State(ctx)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 1, s.SubscriberCount)
	assert.Equal(t, topics, s.Subscribers[0].Topics)

	if !assert.NoError(t, sub.Close()) {
		return
	}

	s, err = pub.State(ctx)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 0, s.SubscriberCount)
}
