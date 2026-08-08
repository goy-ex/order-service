package usecase

import (
	"context"
	"time"

	orderpkg "github.com/goy-ex/order-service/internal/domain/order"
	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
)

type OrderRepo interface {
	InsertWithEvent(ctx context.Context, e *orderpkg.OrderCreated) error
}

// PairStore resolves trading pairs by their key so an order can be validated
// and enriched before being created.
type PairStore interface {
	// ReadByKey returns the pair identified by key, or nil if no such pair is
	// known.
	ReadByKey(key pairpkg.PairKey) *pairpkg.Pair
}

type Clock interface {
	Now() time.Time
}

type ClockFunc func() time.Time

func (f ClockFunc) Now() time.Time {
	return f()
}

type Logger interface {
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
}

type loggerKey int

const LoggerKey loggerKey = 0
