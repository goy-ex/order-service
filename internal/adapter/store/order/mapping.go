package order

import (
	"fmt"

	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	"github.com/goy-ex/order-service/internal/domain"
	orderpkg "github.com/goy-ex/order-service/internal/domain/order"
	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
	"github.com/jackc/pgx/v5/pgtype"
)

func mapOrder(o *orderpkg.Order) pgxsqlc.InsertOrderParams {
	return pgxsqlc.InsertOrderParams{
		ID:           pgtype.UUID{Bytes: o.ID, Valid: true},
		UserID:       pgtype.UUID{Bytes: o.UserID, Valid: true},
		PairBase:     o.Pair.Base,
		PairQuote:    o.Pair.Quote,
		Side:         pgxsqlc.OrderSide(o.Side.String()),
		Price:        int64(o.Price),
		Qty:          int64(o.Qty),
		RemainingQty: int64(o.RemainingQty),
		CreatedAt: pgtype.Timestamptz{
			Time:             o.CreatedAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	}
}

func mapOrderCreated(
	event *orderpkg.OrderCreated,
	encoderFunc OrderCreatedEncoderFunc,
) (pgxsqlc.InsertEventParams, error) {
	payload, err := encoderFunc(event)
	if err != nil {
		var zero pgxsqlc.InsertEventParams

		return zero, fmt.Errorf("failed to encode event: %w", err)
	}

	return pgxsqlc.InsertEventParams{
		ID:            pgtype.UUID{Bytes: event.ID, Valid: true},
		AggregateID:   pgtype.UUID{Bytes: event.Order.ID, Valid: true},
		AggregateType: "order",
		PartitionKey: pgtype.Text{
			String: string(pairpkg.NewPairKey(event.Order.Pair.Base, event.Order.Pair.Quote)),
			Valid:  true,
		},
		EventType: pgxsqlc.EventType(domain.EventTypeOrderCreated.String()),
		Payload:   payload,
		CreatedAt: pgtype.Timestamptz{
			Time:             event.CreatedAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	}, nil
}
