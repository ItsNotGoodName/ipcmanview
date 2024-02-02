package event

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type DahuaDeviceCreated struct {
	Conn models.Conn
}

type DahuaDeviceUpdated struct {
	Conn models.Conn
}

type DahuaDeviceDeleted struct {
	DeviceID int64
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
