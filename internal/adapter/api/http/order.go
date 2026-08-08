package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	goymw "github.com/goy-ex/middleware"
	pkgerrors "github.com/goy-ex/order-service/internal/adapter/api/http/errors"
	"github.com/goy-ex/order-service/internal/usecase"
	"github.com/goy-ex/sentinel"
	"go.uber.org/zap"
)

type orderHandler struct {
	validate           *validator.Validate
	createOrderUseCase usecase.CreateOrderUseCase
}

func NewOrderHandler(
	validate *validator.Validate,
	createOrderUseCase usecase.CreateOrderUseCase,
) *orderHandler {
	return &orderHandler{
		validate:           validate,
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *orderHandler) Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	logger := goymw.GetLogger(r.Context())
	if logger == nil {
		panic(logger)
	}

	var dto CreateOrderDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		err = sentinel.BadRequest(err)
		logger.Info("failed to decode order DTO", zap.Error(err))

		err = pkgerrors.WriteAPIError(w, pkgerrors.ToAPIError(
			fmt.Errorf("%w: %w", pkgerrors.ErrDecodeFail, err),
		))
		if err != nil {
			logger.Warn("failed to write response", zap.Error(err))
		}

		return
	}

	err = h.validate.Struct(dto)
	if err != nil {
		err = sentinel.BadRequest(err)
		logger.Info("failed to validate order DTO", zap.Error(err))

		err = pkgerrors.WriteAPIError(w, pkgerrors.ToAPIError(
			fmt.Errorf("%w: %w", pkgerrors.ErrValidationFail, err),
		))
		if err != nil {
			logger.Warn("failed to write response", zap.Error(err))
		}

		return
	}

	in, err := dto.toInput()
	if err != nil {
		logger.Info("failed to map order DTO", zap.Error(err))

		err = pkgerrors.WriteAPIError(w, pkgerrors.ToAPIError(sentinel.BadRequest(err)))
		if err != nil {
			logger.Warn("failed to write response", zap.Error(err))
		}

		return
	}

	err = h.createOrderUseCase.Invoke(
		context.WithValue(r.Context(), usecase.LoggerKey, newZapLogger(logger)),
		in,
	)
	if err != nil {
		apiErr := pkgerrors.ToAPIError(err)

		if apiErr.Status >= http.StatusInternalServerError {
			logger.Error("failed to create order", zap.Error(err))
		} else {
			logger.Info("failed to create order", zap.Error(err))
		}

		err = pkgerrors.WriteAPIError(w, apiErr)
		if err != nil {
			logger.Warn("failed to write response", zap.Error(err))
		}

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

type zapLogger struct {
	*zap.SugaredLogger
}

func newZapLogger(l *zap.Logger) *zapLogger {
	return &zapLogger{l.WithOptions(zap.AddCallerSkip(1)).Sugar()}
}

func (l *zapLogger) Info(msg string, keysAndValues ...any) {
	l.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...any) {
	l.Warnw(msg, keysAndValues...)
}
