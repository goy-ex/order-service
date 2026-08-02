package usecase

import (
	"context"
	"time"

	"github.com/goy-ex/order-service/internal/domain"
)

type OrderRepo interface {
	InsertWithEvent(ctx context.Context, o *domain.Order, e *domain.OrderCreatedEvent) error
}

// PairStore resolves trading pairs by their key so an order can be validated
// and enriched before being created.
type PairStore interface {
	// ReadByKey returns the pair identified by key, or nil if no such pair is
	// known.
	ReadByKey(key domain.PairKey) *domain.Pair
}

type Event struct {
}

type CreateEvent func(o *domain.Order) *Event

type Clock interface {
	Now() time.Time
}

type ClockFunc func() time.Time

func (f ClockFunc) Now() time.Time {
	return f()
}
