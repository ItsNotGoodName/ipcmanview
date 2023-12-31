package models

type EventDahuaDeviceCreated struct {
	Device DahuaDeviceConn
}

type EventDahuaDeviceUpdated struct {
	Device DahuaDeviceConn
}

type EventDahuaDeviceDeleted struct {
	DeviceID int64
}

type EventDahuaEvent struct {
	Event     DahuaEvent
	EventRule DahuaEventRule
}

type EventDahuaEventWorkerConnecting struct {
	DeviceID int64
}

type EventDahuaEventWorkerConnect struct {
	DeviceID int64
}

type EventDahuaEventWorkerDisconnect struct {
	DeviceID int64
	Error    error
}

type EventDahuaCoaxialStatus struct {
	Channel       int
	CoaxialStatus DahuaCoaxialStatus
}

type EventDahuaQuickScanQueue struct {
	DeviceID int64
}
