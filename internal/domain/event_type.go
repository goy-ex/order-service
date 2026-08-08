package domain

import (
	"strconv"
)

type EventType byte

const (
	eventTypeInvalid      EventType = 0
	EventTypeOrderCreated EventType = 1
)

func (et EventType) IsValid() bool {
	switch et {
	case EventTypeOrderCreated:
		return true
	default:
		return false
	}
}

func (et EventType) String() string {
	switch et {
	case EventTypeOrderCreated:
		return "order_created"
	case eventTypeInvalid:
		return "invalid"
	default:
		return "invalid(" + strconv.Itoa(int(et)) + ")"
	}
}

func EventTypeFromString(s string) EventType {
	switch s {
	case "order_created":
		return EventTypeOrderCreated
	default:
		return eventTypeInvalid
	}
}
