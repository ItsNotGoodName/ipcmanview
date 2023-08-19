package event

import (
	"context"
)

type Bus interface {
	OnBacklog(func(ctx context.Context) error)
	OnDahuaCameraCreated(func(ctx context.Context, evt DahuaCameraCreated) error)
	OnDahuaCameraUpdated(func(ctx context.Context, evt DahuaCameraUpdated) error)
	OnDahuaCameraDeleted(func(ctx context.Context, evt DahuaCameraDeleted) error)
}

type DahuaCameraCreated struct {
	CameraID int64
}

type DahuaCameraUpdated struct {
	CameraID int64
}

type DahuaCameraDeleted struct {
	CameraID int64
}
