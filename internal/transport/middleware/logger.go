package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LogRequest(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().UnixNano())
		reqID := fmt.Sprintf("%016x", rand.Int())[:10]

		ctx := context.WithValue(r.Context(), domains.ReqIDKey{}, reqID)

		middlewareLogger := logger.WithFields(logrus.Fields{
			"request_id":  reqID,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"path":        r.URL.Path,
		})

		// Логгер для передачи в контекст (только request_id)
		contextLogger := logrus.NewEntry(logger).WithField("request_id", reqID) // Важно: создаём новый Entry
		ctx = logctx.WithLogger(ctx, contextLogger)

		middlewareLogger.Info("request started")

		startTime := time.Now()

		defer func() {
			duration := time.Since(startTime)
			middlewareLogger.WithField("duration", duration).Info("request completed")
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
