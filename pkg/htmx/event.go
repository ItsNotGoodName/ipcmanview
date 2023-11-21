package htmx

import (
	"encoding/json"
	"net/http"
)

// Event

type Event struct {
	Key   string
	Value any
}

type Events []Event

func NewEventString(key string) Event {
	return Event{
		Key: key,
	}
}

func NewEvent(key string, value any) Event {
	return Event{
		Key:   key,
		Value: value,
	}
}

func MultipleEvents(events ...Event) Events {
	return events
}

func (e Event) String() string {
	if e.Value == nil {
		return e.Key
	}

	slug := make(map[string]any)
	slug[e.Key] = e.Value

	b, _ := json.Marshal(slug)
	return string(b)
}

func (es Events) String() string {
	slug := make(map[string]any)
	for _, e := range es {
		if e.Value == nil {
			slug[e.Key] = ""
		} else {
			slug[e.Key] = e.Value
		}
	}

	b, _ := json.Marshal(slug)
	return string(b)
}

func (e Event) FromBody() string {
	return string(e.Key) + " from:body"
}

func (e Event) SetTrigger(w http.ResponseWriter) {
	SetTrigger(w, e.String())
}

func (e Event) SetTriggerAfterSettle(w http.ResponseWriter) {
	SetTriggerAfterSettle(w, e.String())
}

func (e Event) SetTriggerAfterSwap(w http.ResponseWriter) {
	SetTriggerAfterSwap(w, e.String())
}

func (es Events) SetTrigger(w http.ResponseWriter) {
	SetTrigger(w, es.String())
}

func (es Events) SetTriggerAfterSettle(w http.ResponseWriter) {
	SetTriggerAfterSettle(w, es.String())
}

func (es Events) SetTriggerAfterSwap(w http.ResponseWriter) {
	SetTriggerAfterSwap(w, es.String())
}
