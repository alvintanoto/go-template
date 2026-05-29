package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type ZapMiddleware struct {
	logger *zap.Logger
}

func NewZapMiddleware(logger *zap.Logger) *ZapMiddleware {
	return &ZapMiddleware{logger: logger}
}

func (m *ZapMiddleware) Execute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		latency := time.Since(start)

		m.logger.Info("HTTP Request Processed",
			zap.String("request_id", middleware.GetReqID(r.Context())),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", ww.Status()),
			zap.Int("bytes_written", ww.BytesWritten()),
			zap.Duration("latency", latency),
			zap.String("remote_ip", r.RemoteAddr),
		)
	})
}
