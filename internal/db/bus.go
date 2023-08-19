package db

import (
	"context"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ event.Bus = (*Bus)(nil)

func NewBusFromPool(pool *pgxpool.Pool) *Bus {
	return NewBus(pool.Config().ConnConfig)
}

func NewBus(config *pgx.ConnConfig) *Bus {
	bus := &Bus{}

	handlers := make(map[string]handlerFunc)
	registerHandlers(bus, handlers)

	bus.listener = listener{
		config:   config,
		handlers: handlers,
		backlog:  bus.backlog,
	}

	return bus
}

func (b *Bus) backlog(ctx context.Context, qes qes.Querier) error {
	for _, v := range b.Backlog {
		if err := v(ctx); err != nil {
			return err
		}
	}

	return nil
}

type Bus struct {
	listener
	Backlog            []func(ctx context.Context) error
	DahuaCameraCreated []func(ctx context.Context, evt event.DahuaCameraCreated) error
	DahuaCameraUpdated []func(ctx context.Context, evt event.DahuaCameraUpdated) error
	DahuaCameraDeleted []func(ctx context.Context, evt event.DahuaCameraDeleted) error
}

func (b *Bus) OnBacklog(fn func(ctx context.Context) error) {
	b.Backlog = append(b.Backlog, fn)
}

func (b *Bus) OnDahuaCameraDeleted(fn func(ctx context.Context, evt event.DahuaCameraDeleted) error) {
	b.DahuaCameraDeleted = append(b.DahuaCameraDeleted, fn)
}

func (b *Bus) OnDahuaCameraUpdated(fn func(ctx context.Context, evt event.DahuaCameraUpdated) error) {
	b.DahuaCameraUpdated = append(b.DahuaCameraUpdated, fn)
}

func (b *Bus) OnDahuaCameraCreated(fn func(ctx context.Context, evt event.DahuaCameraCreated) error) {
	b.DahuaCameraCreated = append(b.DahuaCameraCreated, fn)
}

func registerHandlers(b *Bus, h map[string]handlerFunc) {
	h["dahua.cameras:created"] = func(ctx context.Context, notification *pgconn.Notification, qes qes.Querier) error {
		id, err := strconv.ParseInt(notification.Payload, 10, 64)
		if err != nil {
			return err
		}

		evt := event.DahuaCameraCreated{CameraID: id}
		for _, v := range b.DahuaCameraCreated {
			v(ctx, evt)
		}

		return nil
	}

	h["dahua.cameras:updated"] = func(ctx context.Context, notification *pgconn.Notification, qes qes.Querier) error {
		id, err := strconv.ParseInt(notification.Payload, 10, 64)
		if err != nil {
			return err
		}

		evt := event.DahuaCameraUpdated{CameraID: id}
		for _, v := range b.DahuaCameraUpdated {
			v(ctx, evt)
		}

		return nil
	}

	h["dahua.cameras:deleted"] = func(ctx context.Context, notification *pgconn.Notification, qes qes.Querier) error {
		id, err := strconv.ParseInt(notification.Payload, 10, 64)
		if err != nil {
			return err
		}

		evt := event.DahuaCameraDeleted{CameraID: id}
		for _, v := range b.DahuaCameraDeleted {
			v(ctx, evt)
		}

		return nil
	}
}
