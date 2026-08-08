package event

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/goy-ex/order-service/internal/adapter/event"
	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXEventRepo struct {
	pool *pgxpool.Pool
	q    *pgxsqlc.Queries
}

func NewPGXEventRepo(pool *pgxpool.Pool, q *pgxsqlc.Queries) *PGXEventRepo {
	return &PGXEventRepo{
		pool: pool,
		q:    q,
	}
}

func (r *PGXEventRepo) SelectOrderEventBatch(
	ctx context.Context,
	batchSize int,
) (batch []*event.Event, err error) {
	pairs, err := r.q.SelectPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select pairs: %w", err)
	}

	if len(pairs) == 0 {
		return nil, ErrNoPairs
	}

	rand.Shuffle(len(pairs), func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}

	txq := r.q.WithTx(tx)

	var (
		locked     bool
		lockedPair pgxsqlc.SelectPairsRow
	)

	for _, p := range pairs {
		locked, err = r.q.TryLockPair(ctx, string(pairpkg.NewPairKey(p.PairBase, p.PairQuote)))
		if err != nil {
			return nil, fmt.Errorf("failed to lock pair: %w", err)
		}
		if locked {
			lockedPair = p

			break
		}
	}

	if !locked {
		return nil, ErrNoUnlockedPairs
	}

	batchRows, err := txq.SelectOrderEventBatch(ctx, pgxsqlc.SelectOrderEventBatchParams{
		PairKey: pgtype.Text{
			String: string(pairpkg.NewPairKey(lockedPair.PairBase, lockedPair.PairQuote)),
			Valid:  true,
		},
		BatchSize: int32(batchSize),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to select batch: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
			err = errors.Join(err, rollbackErr)
		}
	}()

	batch, err = mapSelectOrderEventBatchRows(batchRows)
	if err != nil {
		return nil, fmt.Errorf("failed to map batch: %w", err)
	}

	return batch, nil
}

func (r *PGXEventRepo) ResolveBatch(ctx context.Context, batch []*event.Event) error {
	err := r.q.ResolveEventBatch(ctx, mapBatchToUUIDs(batch))
	if err != nil {
		return fmt.Errorf("failed to resolve batch: %w", err)
	}

	return nil
}
