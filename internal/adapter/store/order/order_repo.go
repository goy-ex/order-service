package order

import (
	"context"
	"fmt"

	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxOrderRepo struct {
	pool *pgxpool.Pool
	q    *pgxsqlc.Queries
}

func NewPGXOrderRepo(pool *pgxpool.Pool, q *pgxsqlc.Queries) *pgxOrderRepo {
	return &pgxOrderRepo{
		pool: pool,
		q:    q,
	}
}

func (r *pgxOrderRepo) InsertWithEvent(
	ctx context.Context,
	order *domain.Order,
	event *domain.OrderCreatedEvent,
) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	txq := r.q.WithTx(tx)

	err = txq.InsertOrder(ctx, mapOrder(order))
	if err != nil {
		return fmt.Errorf("sqlc error: %w", err)
	}

	err = txq.InsertEvent(ctx, mapEvent(event))
	if err != nil {
		return fmt.Errorf("sqlc error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
