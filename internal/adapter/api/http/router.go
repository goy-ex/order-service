package http

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	goymw "github.com/goy-ex/middleware"
	mw "github.com/goy-ex/order-service/internal/pkg/middleware"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, orderHandler *orderHandler) *chi.Mux {
	router := chi.NewMux()

	router.Use(
		chimw.RequestID,
		goymw.Recoverer(mw.StandardLogger(logger)),
		goymw.RequestLogger(mw.StandardLogger(logger), "/healthz"),
		chimw.Heartbeat("/healthz"),
	)

	router.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/orders", func(orders chi.Router) {
			orders.Get("/", nil)
			orders.Post("/", orderHandler.Post)
		})
	})

	return router
}
