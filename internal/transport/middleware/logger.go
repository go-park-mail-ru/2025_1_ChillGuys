package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	logger *logrus.Logger
}

func NewLoggerMiddleware(logger *logrus.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (m *LoggerMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := fmt.Sprintf("%016x", rand.Int())[:10]

		ctx := context.WithValue(r.Context(), domains.ReqIDKey, reqId)

		requestLogger := m.logger.WithFields(logrus.Fields{
			"request_id":  reqId,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"path":        r.URL.Path,
		})

		ctx = logctx.WithLogger(ctx, requestLogger)
		r = r.WithContext(ctx)

		requestLogger.Info("request started")

		startTime := time.Now()

		defer func() {
			duration := time.Since(startTime)

			requestLogger.WithField("duration", duration).Info("request completed")
		}()

		next.ServeHTTP(w, r)
		requestLogger.Info("request completed")
	})
}
