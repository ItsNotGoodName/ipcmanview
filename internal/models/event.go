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
	Event DahuaEvent
}

// TODO: these should be generated

func (EventDahuaCameraDeleted) EventName() string {
	return "EventDahuaCameraDeleted"
}

func (EventDahuaCameraCreated) EventName() string {
	return "EventDahuaCameraCreated"
}

func (EventDahuaCameraUpdated) EventName() string {
	return "EventDahuaCameraUpdated"
}

func (e EventDahuaCameraEvent) EventName() string {
	return "EventDahuaCameraEvent"
}
