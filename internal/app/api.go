package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	httpapi "github.com/goy-ex/order-service/internal/adapter/api/http"
	"github.com/goy-ex/order-service/internal/adapter/store"
	orderstore "github.com/goy-ex/order-service/internal/adapter/store/order"
	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	"github.com/goy-ex/order-service/internal/domain"
	pairpkg "github.com/goy-ex/order-service/internal/domain/pair"
	configpkg "github.com/goy-ex/order-service/internal/pkg/config"
	"github.com/goy-ex/order-service/internal/pkg/event"
	pkghttp "github.com/goy-ex/order-service/internal/pkg/http"
	"github.com/goy-ex/order-service/internal/pkg/kvdecoder"
	"github.com/goy-ex/order-service/internal/pkg/loader"
	"github.com/goy-ex/order-service/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type api struct {
	pool        *pgxpool.Pool
	server      *http.Server
	config      *configpkg.APIConfig
	pairWorker  *kvPoller[pairpkg.PairKey, *pairpkg.Pair]
	assetWorker *kvPoller[string, *domain.Asset]
}

func NewAPI(ctx context.Context, config *configpkg.APIConfig, logger *zap.Logger) (*api, error) {
	validate := validator.New()

	pairStore := store.NewPointerKVStore[pairpkg.PairKey, *pairpkg.Pair]()
	pairLoader := loader.NewLoader(
		loader.NewFileFetcher(config.Pairs.Source),
		loader.DecoderFunc[map[pairpkg.PairKey]*pairpkg.Pair](kvdecoder.DecodePairsJSON),
	)

	pairs, err := pairLoader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load pairs: %w", err)
	}

	pairStore.Set(pairs)

	pairWorker := newKVPoller(
		logger,
		pairLoader,
		pairStore,
		time.Duration(config.Pairs.PollingRate)*time.Second,
	)

	assetStore := store.NewPointerKVStore[string, *domain.Asset]()
	assetLoader := loader.NewLoader(
		loader.NewFileFetcher(config.Assets.Source),
		loader.DecoderFunc[map[string]*domain.Asset](kvdecoder.DecodeAssetsJSON),
	)

	assets, err := assetLoader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load assets: %w", err)
	}

	assetStore.Set(assets)

	assetWorker := newKVPoller(
		logger,
		assetLoader,
		assetStore,
		time.Duration(config.Assets.PollingRate)*time.Second,
	)

	pool, err := pgxpool.New(ctx, config.DBConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	queries := pgxsqlc.New(pool)
	orderRepo := orderstore.NewPGXOrderRepo(pool, queries, event.EncodeOrderCreatedToJSON)

	createOrderUseCase := usecase.NewCreateOrderUseCase(
		pairStore,
		orderRepo,
		usecase.ClockFunc(time.Now),
	)

	orderHandler := httpapi.NewOrderHandler(validate, createOrderUseCase)
	router := httpapi.NewRouter(logger, orderHandler)
	server := pkghttp.NewServer(&pkghttp.ServerConfig{
		Addr:              config.Server.Addr,
		ReadHeaderTimeout: config.Server.ReadHeaderTimeout,
		ReadTimeout:       config.Server.ReadTimeout,
		WriteTimeout:      config.Server.WriteTimeout,
	}, router)

	return &api{
		pool:        pool,
		server:      server,
		config:      config,
		pairWorker:  pairWorker,
		assetWorker: assetWorker,
	}, nil
}

func (api *api) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		api.pairWorker.run(ctx)
	}()

	go func() {
		api.assetWorker.run(ctx)
	}()

	go func() {
		err := api.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(api.config.ShutdownTimeout)*time.Second,
		)
		defer cancel()

		//nolint:contextcheck // shutdown context
		err := api.server.Shutdown(shutdownCtx)
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		return nil
	case err := <-errCh:
		return err
	}
}

func (api *api) Close() {
	api.pool.Close()
}
