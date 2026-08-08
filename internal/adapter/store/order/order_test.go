package order_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	orderstore "github.com/goy-ex/order-service/internal/adapter/store/order"
	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	orderpkg "github.com/goy-ex/order-service/internal/domain/order"
	"github.com/goy-ex/order-service/internal/pkg/event"
	"github.com/goy-ex/order-service/internal/pkg/fixture"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18",
		postgres.WithInitScripts(
			"./testdata/schema.sql",
		),
		postgres.WithDatabase("order_service_db"),
		postgres.WithUsername("order_service_admin"),
		postgres.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	connString, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	pool, err = pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("failed to create pool")
	}

	code := m.Run()
	os.Exit(code)
}

func TestOrderRepoInsertWithEvent_Valid(t *testing.T) {
	t.Parallel()

	tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})
	require.NoError(t, err)

	queries := pgxsqlc.New(tx)

	orderRepo := orderstore.NewPGXOrderRepo(pool, queries, event.EncodeOrderCreatedToJSON)

	id, err := uuid.NewV7()
	require.NoError(t, err)

	order := fixture.NewValidOrder()
	event := orderpkg.NewOrderCreated(id, &order)

	err = orderRepo.InsertWithEvent(context.Background(), event)
	require.NoError(t, err)

	tx.Rollback(context.Background())
}

// func TestOrderRepoInsertWithEvent_DuplicateOrder(t *testing.T) {
// 	t.Parallel()

// 	tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})
// 	require.NoError(t, err)

// 	queries := pgxsqlc.New(tx)

// 	orderRepo := orderstore.NewPGXOrderRepo(pool, queries, encoding.EncodeOrderCreatedToJSON)

// 	order := fixture.NewValidOrder()
// 	event := domain.NewOrderCreated(&order)
// 	err = orderRepo.InsertWithEvent(context.Background(), event)
// 	require.NoError(t, err)

// 	err = orderRepo.InsertWithEvent(context.Background(), event)
// 	require.Error(t, err)

// 	tx.Rollback(context.Background())
// }

// func TestOrderRepoInsertWithEvent_DuplicateEvent(t *testing.T) {
// 	t.Parallel()

// 	tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})
// 	require.NoError(t, err)

// 	queries := pgxsqlc.New(tx)

// 	orderRepo := orderstore.NewPGXOrderRepo(pool, queries, encoding.EncodeOrderCreatedToJSON)

// 	order := fixture.NewValidOrder()
// 	event := domain.NewOrderCreated(&order)
// 	err = orderRepo.InsertWithEvent(context.Background(), event)
// 	require.NoError(t, err)

// 	order2 := fixture.NewValidOrder()
// 	err = orderRepo.InsertWithEvent(context.Background(), event)
// 	require.Error(t, err)

// 	tx.Rollback(context.Background())
// }
