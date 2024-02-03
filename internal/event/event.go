package event

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type DahuaDeviceChanged struct {
	DeviceID int64
	Created  bool
	Updated  bool
	Deleted  bool
	Disabled bool
	Enabled  bool
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
