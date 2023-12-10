package models

type EventDahuaCameraCreated struct {
	Camera DahuaCamera
}

func (EventDahuaCameraCreated) EventName() string {
	return "EventDahuaCameraCreated"
}

type EventDahuaCameraUpdated struct {
	Camera DahuaCamera
}

func (EventDahuaCameraUpdated) EventName() string {
	return "EventDahuaCameraUpdated"
}

type EventDahuaCameraDeleted struct {
	CameraID int64
}

func (EventDahuaCameraDeleted) EventName() string {
	return "EventDahuaCameraDeleted"
}

type EventDahuaCameraEvent struct {
	Event DahuaEvent `json:"event"`
}

func (e EventDahuaCameraEvent) EventName() string {
	return "EventDahuaCameraEvent"
}
