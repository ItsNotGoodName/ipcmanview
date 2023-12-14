package models

type EventDahuaCameraCreated struct {
	Camera DahuaConn
}

type EventDahuaCameraUpdated struct {
	Camera DahuaConn
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

// TODO: these should be generated

func (EventDahuaCameraDeleted) EventTopic() string {
	return "EventDahuaCameraDeleted"
}

func (EventDahuaCameraCreated) EventTopic() string {
	return "EventDahuaCameraCreated"
}

func (EventDahuaCameraUpdated) EventTopic() string {
	return "EventDahuaCameraUpdated"
}

func (e EventDahuaCameraEvent) EventTopic() string {
	return "EventDahuaCameraEvent"
}
