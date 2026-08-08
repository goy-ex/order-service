package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/goy-ex/order-service/internal/domain"
	orderpkg "github.com/goy-ex/order-service/internal/domain/order"
	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
	"github.com/goy-ex/sentinel"
)

// CreateOrderUseCase validates an incoming order against its trading pair and
// hands it off to be matched.
type CreateOrderUseCase interface {
	// Invoke runs the create-order flow for in, returning an error if the pair
	// is unknown, the order is invalid, or it could not be published.
	Invoke(ctx context.Context, in *CreateOrderInput) error
}

type createOrderUseCase struct {
	pairs  PairStore
	orders OrderRepo
	clock  Clock
}

// NewCreateOrderUseCase returns a CreateOrderUseCase that resolves pairs via
// pairReader and publishes created orders via orderProducer.
func NewCreateOrderUseCase(
	pairs PairStore,
	orders OrderRepo,
	clock Clock,
) *createOrderUseCase {
	return &createOrderUseCase{
		pairs:  pairs,
		orders: orders,
		clock:  clock,
	}
}

// CreateOrderInput carries the parameters of an incoming order exactly as
// received from the transport layer, before pair resolution and domain
// validation.
type CreateOrderInput struct {
	Base         string
	Quote        string
	Side         string
	Price        int64
	Qty          int64
	RemainingQty int64
	UserID       uuid.UUID
}

func (uc *createOrderUseCase) Invoke(ctx context.Context, in *CreateOrderInput) error {
	logger, ok := ctx.Value(LoggerKey).(Logger)
	if !ok {
		return sentinel.InvariantViolation(errors.New("failed to get logger"))
	}

	logger.Info("success getting logger")

	if in == nil {
		return sentinel.BadRequest(ErrNilInput)
	}

	pairKey := pairpkg.NewPairKey(in.Base, in.Quote)

	pair := uc.pairs.ReadByKey(pairKey)
	if pair == nil {
		return sentinel.BadRequest(&PairNotFoundError{PairKey: string(pairKey)})
	}

	orderID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate order's UUID (v7): %w", err)
	}

	order, err := orderpkg.NewOrder(
		orderID,
		in.UserID,
		pair,
		domain.SideFromString(in.Side),
		int(in.Price),
		int(in.Qty),
		int(in.RemainingQty),
		uc.clock.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	event := orderpkg.NewOrderCreated(id, order)

	err = uc.orders.InsertWithEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}
