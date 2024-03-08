package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
)

const levelDefault = models.DahuaPermissionLevel_User
const levelEmail = models.DahuaPermissionLevel_User

func Level(ctx context.Context, deviceID int64, level models.DahuaPermissionLevel) (bool, error) {
	actor := core.UseActor(ctx)
	if actor.Admin {
		return true, nil
	}

	dbLevel, err := app.DB.C().DahuaGetPermissionLevel(ctx, repo.DahuaGetPermissionLevelParams{
		DeviceID: deviceID,
		UserID:   core.Int64ToNullInt64(actor.UserID),
	})
	if err != nil {
		return false, err
	}
	return dbLevel >= level, nil
}

func PubSubMiddleware(ctx context.Context) pubsub.MiddlewareFunc {
	actor := core.UseActor(ctx)

	skip := func(ctx context.Context, deviceID int64, level models.DahuaPermissionLevel) bool {
		dbLevel, err := app.DB.C().DahuaGetPermissionLevel(ctx, repo.DahuaGetPermissionLevelParams{
			DeviceID: deviceID,
			UserID:   core.Int64ToNullInt64(actor.UserID),
		})
		if err != nil {
			return true
		}
		return dbLevel <= level
	}

	return func(next pubsub.HandleFunc) pubsub.HandleFunc {
		return func(ctx context.Context, event pubsub.Event) error {
			if actor.Admin {
				return next(ctx, event)
			}

			switch e := event.(type) {
			case bus.DahuaEvent:
				if skip(ctx, e.Event.DeviceID, levelDefault) {
					return nil
				}
			case bus.DahuaFileCreated:
				if skip(ctx, e.DeviceID, levelDefault) {
					return nil
				}
			case bus.DahuaFileCursorUpdated:
				if skip(ctx, e.Cursor.DeviceID, levelDefault) {
					return nil
				}
			case bus.DahuaEmailCreated:
				if skip(ctx, e.DeviceID, levelEmail) {
					return nil
				}
			case bus.UserSecurityUpdated:
				if e.UserID != actor.UserID {
					return nil
				}
			default:
				return nil
			}

			return next(ctx, event)
		}
	}
}
