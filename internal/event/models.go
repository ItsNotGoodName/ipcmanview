package event

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type EventQueued struct {
}

type Event struct {
	Event repo.Event
}

type DahuaEvent struct {
	Event     repo.DahuaEvent
	EventRule repo.DahuaEventRule
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
