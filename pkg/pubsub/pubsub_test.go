package pubsub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func noopEventHandler(ctx context.Context, event Event) error { return nil }

func TestPub(t *testing.T) {
	pub := NewPub()

	sub, err := pub.
		Subscribe().
		Function(noopEventHandler)
	if !assert.NoError(t, err) {
		return
	}
	defer sub.Close()

	s, err := pub.State()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 1, s.SubscriberCount)

	sub.Close()

	s, err = pub.State()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 0, s.SubscriberCount)
}
