package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain/order"
)

type orderCreatedPayload struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	PairBase     string    `json:"pairBase"`
	PairQuote    string    `json:"pairQuote"`
	Side         string    `json:"side"`
	Price        int       `json:"price"`
	Qty          int       `json:"qty"`
	RemainingQty int       `json:"remainingQty"`
	CreatedAt    time.Time `json:"createdAt"`
}

func EncodeOrderCreatedToJSON(event *order.OrderCreated) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	enc := json.NewEncoder(buf)

	payload := &orderCreatedPayload{
		ID:           event.Order.ID,
		UserID:       event.Order.UserID,
		PairBase:     event.Order.Pair.Base,
		PairQuote:    event.Order.Pair.Quote,
		Side:         event.Order.Side.String(),
		Price:        event.Order.Price,
		Qty:          event.Order.Qty,
		RemainingQty: event.Order.RemainingQty,
		CreatedAt:    event.Order.CreatedAt,
	}

	err := enc.Encode(payload)

	if err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}

	return buf.Bytes(), nil
}
