package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goy-ex/order-service/internal/adapter/event"
	eventstore "github.com/goy-ex/order-service/internal/adapter/store/event"
	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	pkgconfig "github.com/goy-ex/order-service/internal/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type OutboxWorker struct {
	logger   *zap.Logger
	config   *pkgconfig.OutboxWorkerConfig
	pool     *pgxpool.Pool
	events   *eventstore.PGXEventRepo
	producer *event.KafkaOrderProducer
}

func NewOutboxWorker(
	ctx context.Context,
	config *pkgconfig.OutboxWorkerConfig,
	logger *zap.Logger,
) (*OutboxWorker, error) {
	pool, err := pgxpool.New(ctx, config.DBConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	queries := pgxsqlc.New(pool)
	events := eventstore.NewPGXEventRepo(pool, queries)
	producer, err := event.NewKafkaOrderProducer(config.OrdersTopic, strings.Split(config.KafkaBrokers, ",")...)
	if err != nil {
		return nil, fmt.Errorf("failed to create order producer: %w", err)
	}

	return &OutboxWorker{
		logger:   logger,
		config:   config,
		pool:     pool,
		events:   events,
		producer: producer,
	}, nil
}

func (ow *OutboxWorker) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(ow.config.PollingRate) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ow.processBatch(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (ow *OutboxWorker) Close() {
	ow.pool.Close()
}

func (ow *OutboxWorker) processBatch(ctx context.Context) {
	batch, err := ow.events.SelectOrderEventBatch(ctx, ow.config.BatchSize)
	if err != nil {
		var log func(msg string, fields ...zap.Field)
		if errors.Is(err, eventstore.ErrNoPairs) || errors.Is(err, eventstore.ErrNoUnlockedPairs) {
			log = ow.logger.Info
		} else {
			log = ow.logger.Error
		}

		log("failed to process batch", zap.Error(err))

		return
	}

	ow.logger.Debug("batch selected", zap.Any("batch", batch))

	if len(batch) == 0 {
		ow.logger.Debug("empty batch")

		return
	}

	err = ow.producer.Produce(ctx, batch)
	if err != nil {
		ow.logger.Error("failed to produce batch", zap.Error(err))

		return
	}

	ow.logger.Debug("batch produced")

	err = ow.events.ResolveBatch(ctx, batch)
	if err != nil {
		ow.logger.Error("failed to resolve batch", zap.Error(err))

		return
	}

	ow.logger.Debug("batch resolved")
}
