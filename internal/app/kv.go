package app

import (
	"context"
	"time"

	pkgstore "github.com/goy-ex/order-service/internal/adapter/store"
	pkgloader "github.com/goy-ex/order-service/internal/pkg/loader"
	"go.uber.org/zap"
)

type kvPoller[K comparable, V any] struct {
	logger *zap.Logger
	loader *pkgloader.Loader[map[K]V]
	store  *pkgstore.PointerKVStore[K, V]
	rate   time.Duration
}

func newKVPoller[K comparable, V any](
	logger *zap.Logger,
	loader *pkgloader.Loader[map[K]V],
	store *pkgstore.PointerKVStore[K, V],
	rate time.Duration,
) *kvPoller[K, V] {
	return &kvPoller[K, V]{
		logger: logger,
		loader: loader,
		store:  store,
		rate:   rate,
	}
}

func (p kvPoller[K, V]) run(ctx context.Context) {
	ticker := time.NewTicker(p.rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			kv, err := p.loader.Load()
			if err != nil {
				p.logger.Warn("failed to load kv", zap.Error(err))

				continue
			}

			p.store.Set(kv)
		case <-ctx.Done():
			return
		}
	}
}
