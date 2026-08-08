package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreated struct {
	CreatedAt time.Time
	Order     *Order
	ID        uuid.UUID
}

func NewOrderCreated(id uuid.UUID, o *Order) *OrderCreated {
	return &OrderCreated{
		ID:        id,
		Order:     o,
		CreatedAt: o.CreatedAt,
	}
}
