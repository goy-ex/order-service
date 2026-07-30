package domain

import "strconv"

type EventType byte

const (
	EventTypeOrderCreated EventType = iota
)

func (t EventType) String() string {
	switch t {
	case EventTypeOrderCreated:
		return "order created"
	default:
		return "invalid(" + strconv.Itoa(int(t)) + ")"
	}
}

type Event interface {
	GetEventType() EventType
}
