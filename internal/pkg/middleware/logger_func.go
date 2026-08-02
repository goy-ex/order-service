package middleware

import (
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
	goymw "github.com/goy-ex/middleware"
	"go.uber.org/zap"
)

func StandardLogger(logger *zap.Logger) goymw.LoggerFunc {
	return func(r *http.Request) *zap.Logger {
		reqID := chimw.GetReqID(r.Context())
		if reqID == "" {
			reqID = "unknown"
		}

		return logger.With(
			zap.String("reqID", reqID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remoteAddr", r.RemoteAddr),
		)
	}
}
