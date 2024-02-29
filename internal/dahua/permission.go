package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pubsub"
)

const levelDefault = models.DahuaPermissionLevelUser
const levelEmail = models.DahuaPermissionLevelUser

func Level(ctx context.Context, db sqlite.DB, deviceID int64, level models.DahuaPermissionLevel) (bool, error) {
	actor := core.UseActor(ctx)
	if actor.Admin {
		return true, nil
	}

	dbLevel, err := db.C().DahuaGetPermissionLevel(ctx, repo.DahuaGetPermissionLevelParams{
		DeviceID: deviceID,
		UserID:   core.Int64ToNullInt64(actor.UserID),
	})
	if err != nil {
		return false, err
	}
	return dbLevel >= level, nil
}

func PubSubMiddleware(ctx context.Context, db sqlite.DB) pubsub.MiddlewareFunc {
	actor := core.UseActor(ctx)

	skip := func(ctx context.Context, deviceID int64, level models.DahuaPermissionLevel) bool {
		dbLevel, err := db.C().DahuaGetPermissionLevel(ctx, repo.DahuaGetPermissionLevelParams{
			DeviceID: deviceID,
			UserID:   core.Int64ToNullInt64(actor.UserID),
		})
		if err != nil {
			return true
		}
		return dbLevel <= level
	}

	return func(next pubsub.HandleFunc) pubsub.HandleFunc {
		return func(ctx context.Context, evt pubsub.Event) error {
			if actor.Admin {
				return next(ctx, evt)
			}

			switch e := evt.(type) {
			case event.DahuaEvent:
				if skip(ctx, e.Event.DeviceID, levelDefault) {
					return nil
				}
			case event.DahuaFileCreated:
				if skip(ctx, e.DeviceID, levelDefault) {
					return nil
				}
			case event.DahuaFileCursorUpdated:
				if skip(ctx, e.Cursor.DeviceID, levelDefault) {
					return nil
				}
			case event.DahuaEmailCreated:
				if skip(ctx, e.DeviceID, levelEmail) {
					return nil
				}
			case event.UserSecurityUpdated:
				if e.UserID != actor.UserID {
					return nil
				}
			default:
				return nil
			}

			return next(ctx, evt)
		}
	}
}
