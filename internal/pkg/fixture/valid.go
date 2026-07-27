package fixture

import (
	"time"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain"
)

func NewValidPair() domain.Pair {
	pair, err := domain.NewPair("BTC", "USDT", 1, 1, 1)
	if err != nil {
		panic(err)
	}

	return *pair
}

func NewValidOrder() domain.Order {
	pair := NewValidPair()

	order, err := domain.NewOrder(
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
