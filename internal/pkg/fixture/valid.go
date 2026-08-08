package fixture

import (
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/goy-ex/order-service/internal/domain/order"
	"github.com/goy-ex/order-service/internal/domain/pair"
)

func NewValidPair() pair.Pair {
	pair, err := pair.NewPair("BTC", "USDT", 1, 1, 1)
	if err != nil {
		panic(err)
	}

	return *pair
}

func NewValidOrder() order.Order {
	pair := NewValidPair()

	order, err := order.NewOrder(
		uuid.New(),
		uuid.New(),
		&pair,
		domain.SideAsk,
		1,
		1,
		1,
		time.Now(),
	)
	if err != nil {
		panic(err)
	}

	return *order
}
