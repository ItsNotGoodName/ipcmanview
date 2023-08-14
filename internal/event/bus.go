package event

import (
	"context"
)

type Bus interface {
	OnBacklog(func(ctx context.Context) error)
	OnDahuaCameraUpdated(func(ctx context.Context, evt DahuaCameraUpdated) error)
	OnDahuaCameraDeleted(func(ctx context.Context, evt DahuaCameraDeleted) error)
}

type DahuaCameraUpdated struct {
	CameraID int64
}

type DahuaCameraDeleted struct {
	CameraID int64
}
