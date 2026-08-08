package order

import (
	"context"
	"fmt"

	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	orderpkg "github.com/goy-ex/order-service/internal/domain/order"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderCreatedEncoderFunc func(*orderpkg.OrderCreated) ([]byte, error)

type pgxOrderRepo struct {
	pool                    *pgxpool.Pool
	q                       *pgxsqlc.Queries
	orderCreatedEncoderFunc OrderCreatedEncoderFunc
}

func NewPGXOrderRepo(
	pool *pgxpool.Pool,
	q *pgxsqlc.Queries,
	orderCreatedEncoderFunc OrderCreatedEncoderFunc,
) *pgxOrderRepo {
	return &pgxOrderRepo{
		pool:                    pool,
		q:                       q,
		orderCreatedEncoderFunc: orderCreatedEncoderFunc,
	}
}

func (r *pgxOrderRepo) InsertWithEvent(
	ctx context.Context,
	event *orderpkg.OrderCreated,
) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	txq := r.q.WithTx(tx)

	err = txq.InsertOrder(ctx, mapOrder(event.Order))
	if err != nil {
		return fmt.Errorf("sqlc error: %w", err)
	}

	arg, err := mapOrderCreated(event, r.orderCreatedEncoderFunc)
	if err != nil {
		return fmt.Errorf("failed to map order created event: %w", err)
	}

	err = txq.InsertEvent(ctx, arg)
	if err != nil {
		return fmt.Errorf("sqlc error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
