package event

import (
	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain"
)

type Event struct {
	ID            uuid.UUID
	AggregateID   uuid.UUID
	AggregateType string
	PartitionKey  string
	EventType     domain.EventType
	Payload       []byte
}
