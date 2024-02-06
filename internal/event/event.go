package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

const (
	ActionDahuaDeviceCreated models.EventAction = "dahua-device:created"
	ActionDahuaDeviceUpdated models.EventAction = "dahua-device:updated"
	ActionDahuaDeviceDeleted models.EventAction = "dahua-device:deleted"
)

func CreateEvent(ctx context.Context, db sqlite.DBTx, action models.EventAction, data any) (int64, error) {
	actor := core.UseActor(ctx)
	b, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return db.CreateEvent(ctx, repo.CreateEventParams{
		Action: action,
		Data:   types.NewJSON(b),
		UserID: sql.NullInt64{
			Int64: actor.UserID,
			Valid: actor.Type == core.ActorTypeUser,
		},
		Actor:     actor.Type,
		CreatedAt: types.NewTime(time.Now()),
	})
}

func UseDataDahuaDevice(evt repo.Event) int64 {
	var deviceID int64
	err := json.Unmarshal(evt.Data.RawMessage, &deviceID)
	if err != nil {
		return 0
	}
	return deviceID
}

type EventQueued struct {
}

type Event struct {
	Event repo.Event
}

type DahuaEvent struct {
	DeviceName string
	Event      repo.DahuaEvent
	EventRule  repo.DahuaEventRule
}

type DahuaEventWorkerConnecting struct {
	DeviceID int64
}

type DahuaEventWorkerConnect struct {
	DeviceID int64
}

type DahuaEventWorkerDisconnect struct {
	DeviceID int64
	Error    error
}

type DahuaCoaxialStatus struct {
	DeviceID      int64
	Channel       int
	CoaxialStatus models.DahuaCoaxialStatus
}
