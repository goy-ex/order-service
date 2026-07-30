package order

import (
	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func mapOrder(o *domain.Order) pgxsqlc.InsertOrderParams {
	return pgxsqlc.InsertOrderParams{
		ID:        pgtype.UUID{Bytes: o.ID, Valid: true},
		UserID:    pgtype.UUID{Bytes: o.UserID, Valid: true},
		PairBase:  o.Pair.Base,
		PairQuote: o.Pair.Quote,
		Side:      pgxsqlc.OrderSide(o.Side.String()),
		Price:     int64(o.Price),
		Amount:    int64(o.Amount),
		Remaining: int64(o.Remaining),
		CreatedAt: pgtype.Timestamptz{
			Time:             o.CreatedAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
	}
}

func mapEvent(event *domain.OrderCreatedEvent) pgxsqlc.InsertEventParams {
	return pgxsqlc.InsertEventParams{
		ID:            pgtype.UUID{Bytes: event.ID, Valid: true},
		AggregateID:   pgtype.UUID{Bytes: event.Order.ID, Valid: true},
		AggregateType: "order",
		EventType:     domain.EventTypeOrderCreated.String(),
		CreatedAt: pgtype.Timestamptz{
			Time:             event.CreatedAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		Payload:    []byte("{}"),
		IsCommited: false,
	}
}
