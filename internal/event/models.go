package event

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

type UserSecurityUpdated struct {
	UserID int64
}

type DahuaEvent struct {
	Event     repo.DahuaEvent
	EventRule repo.DahuaEventRule
}

type DahuaDeviceCreated struct {
	DeviceID int64
}

type DahuaDeviceUpdated struct {
	DeviceID int64
}

type DahuaDeviceDeleted struct {
	DeviceID int64
}

type DahuaEmailCreated struct {
	DeviceID int64
	EmailID  int64
}

type DahuaFileCreated struct {
	DeviceID  int64
	TimeRange models.TimeRange
	Count     int64
}

type DahuaFileCursorUpdated struct {
	Cursor repo.DahuaFileCursor
}

type DahuaWorkerConnecting struct {
	DeviceID int64
	Type     models.DahuaWorkerType
}

type DahuaWorkerConnected struct {
	DeviceID int64
	Type     models.DahuaWorkerType
}

type DahuaWorkerDisconnected struct {
	DeviceID int64
	Type     models.DahuaWorkerType
}

type DahuaCoaxialStatus struct {
	DeviceID      int64
	Channel       int
	CoaxialStatus models.DahuaCoaxialStatus
}
