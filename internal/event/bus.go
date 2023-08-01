package event

import (
	"fmt"
	"strconv"

	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type Bus struct {
	DahuaCameraUpdated []func(context db.Context, evt DahuaCameraUpdated)
	DahuaCameraDeleted []func(context db.Context, evt DahuaCameraDeleted)
}

type DahuaCameraUpdated struct {
	IDS []int64
}

type DahuaCameraDeleted struct {
	IDS []int64
}

var (
	dahuaCamerasUpdated = "dahua_cameras:updated"
	dahuaCamerasDeleted = "dahua_cameras:deleted"
)

var channels = []string{
	dahuaCamerasUpdated,
	dahuaCamerasDeleted,
}

func (b *Bus) Handle(dbCtx db.Context, notification *pgconn.Notification) {
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
			v(dbCtx, evt)
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
			v(dbCtx, evt)
		}
	}
}

func NewBus() *Bus {
	return &Bus{}
}
