package event

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

const (
	ActionDahuaDeviceCreated models.EventAction = "dahua-device:created"
	ActionDahuaDeviceUpdated models.EventAction = "dahua-device:updated"
	ActionDahuaDeviceDeleted models.EventAction = "dahua-device:deleted"
)

type EventQueued struct {
}

type EventCreated struct {
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
