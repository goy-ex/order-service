package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/sentinel"
)

//nolint:govet // entity
type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Pair      *Pair
	Side      Side
	Price     int
	Amount    int
	Remaining int
	CreatedAt time.Time
}

func NewOrder(
	id uuid.UUID,
	userID uuid.UUID,
	pair *Pair,
	side Side,
	price int,
	amount int,
	remaining int,
	createdAt time.Time,
) (*Order, error) {
	if price <= 0 {
		return nil, sentinel.BadRequest(&FieldNotPositiveError{FieldName: "Price", Value: price})
	}

	if amount <= 0 {
		return nil, sentinel.BadRequest(&FieldNotPositiveError{FieldName: "Amount", Value: amount})
	}

	if remaining <= 0 {
		return nil, sentinel.BadRequest(&FieldNotPositiveError{FieldName: "Remaining", Value: remaining})
	}

	return &Order{
		ID:        id,
		UserID:    userID,
		Side:      side,
		Pair:      pair,
		Price:     price,
		Amount:    amount,
		Remaining: remaining,
		CreatedAt: createdAt,
	}, nil
}
