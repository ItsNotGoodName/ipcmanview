package event

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type Bus struct {
	Connect            []func(ctx context.Context)
	DahuaCameraUpdated []func(ctx context.Context, evt DahuaCameraUpdated)
	DahuaCameraDeleted []func(ctx context.Context, evt DahuaCameraDeleted)
}

type DahuaCameraUpdated struct {
	IDS []int64
}

type DahuaCameraDeleted struct {
	IDS []int64
}

var (
	dahuaCamerasUpdated = "dahua.cameras:updated"
	dahuaCamerasDeleted = "dahua.cameras:deleted"
)

var channels = []string{
	dahuaCamerasUpdated,
	dahuaCamerasDeleted,
}

func (b *Bus) handle(ctx context.Context, notification *pgconn.Notification) {
	switch notification.Channel {
	case dahuaCamerasDeleted:
		id, err := strconv.ParseInt(notification.Payload, 10, 64)
		if err != nil {
			log.Err(err).Str("payload", notification.Payload).Msg("Invalid payload from notification")
			return
		}

		fmt.Println("Camera deleted ", id)

		evt := DahuaCameraDeleted{IDS: []int64{id}}
		for _, v := range b.DahuaCameraDeleted {
			v(ctx, evt)
		}
	case dahuaCamerasUpdated:
		id, err := strconv.ParseInt(notification.Payload, 10, 64)
		if err != nil {
			log.Err(err).Str("payload", notification.Payload).Msg("Invalid payload from notification")
			return
		}

		fmt.Println("Camera updated ", id)

		evt := DahuaCameraUpdated{IDS: []int64{id}}
		for _, v := range b.DahuaCameraUpdated {
			v(ctx, evt)
		}
	}
}

func (b *Bus) handleConnect(ctx context.Context) {
	for _, v := range b.Connect {
		v(ctx)
	}
}

func NewBus() *Bus {
	return &Bus{}
}
