package models

type EventDahuaCameraCreated struct {
	Camera DahuaCameraConn
}

type EventDahuaCameraUpdated struct {
	Camera DahuaCameraConn
}

type EventDahuaCameraDeleted struct {
	CameraID int64
}

type EventDahuaCameraEvent struct {
	Event     DahuaEvent
	EventRule DahuaEventRule
}

type EventDahuaEventWorkerConnecting struct {
	CameraID int64
}

type EventDahuaEventWorkerConnect struct {
	CameraID int64
}

type EventDahuaEventWorkerDisconnect struct {
	CameraID int64
	Error    error
}

type EventDahuaCoaxialStatus struct {
	Channel       int
	CoaxialStatus DahuaCoaxialStatus
}
