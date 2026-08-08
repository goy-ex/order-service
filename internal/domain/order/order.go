package order

import (
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/goy-ex/order-service/internal/domain/pair"
	"github.com/goy-ex/sentinel"
)

//nolint:govet // entity
type Order struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Pair         *pair.Pair
	Side         domain.Side
	Price        int
	Qty          int
	RemainingQty int
	CreatedAt    time.Time
}

func NewOrder(
	id uuid.UUID,
	userID uuid.UUID,
	pair *pair.Pair,
	side domain.Side,
	price int,
	qty int,
	remainingQty int,
	createdAt time.Time,
) (*Order, error) {
	if price <= 0 {
		return nil, sentinel.BadRequest(&domain.NotPositiveNumberError{FieldName: "Price", Value: price})
	}

	if qty <= 0 {
		return nil, sentinel.BadRequest(&domain.NotPositiveNumberError{FieldName: "Amount", Value: qty})
	}

	if remainingQty <= 0 {
		return nil, sentinel.BadRequest(&domain.NotPositiveNumberError{FieldName: "Remaining", Value: remainingQty})
	}

	if !side.IsValid() {
		return nil, sentinel.BadRequest(&domain.NoSuchOptionError{FieldName: "Side", Value: side.String()})
	}

	order := &Order{
		ID:           id,
		UserID:       userID,
		Side:         side,
		Pair:         pair,
		Price:        price,
		Qty:          qty,
		RemainingQty: remainingQty,
		CreatedAt:    createdAt,
	}

	return order, nil
}
