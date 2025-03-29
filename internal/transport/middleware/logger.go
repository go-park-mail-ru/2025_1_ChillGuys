package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ContextKey string

const ReqIDKey ContextKey = "ReqId"

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
		// Генерируем уникальный ID запроса
		reqId := fmt.Sprintf("%016x", rand.Int())[:10]

		// Добавляем ID в контекст
		ctx := context.WithValue(r.Context(), ReqIDKey, reqId)
		r = r.WithContext(ctx)

		// Создаем логгер с полями для этого запроса
		requestLogger := m.logger.WithFields(logrus.Fields{
			"request_id":  reqId,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"path":        r.URL.Path,
		})

		// Cохраняем логгер в контекст для использования в обработчиках
		ctx = context.WithValue(ctx, ContextKey("logger"), requestLogger)
		r = r.WithContext(ctx)

		requestLogger.Info("request started")

		// Передаем запрос дальше
		next.ServeHTTP(w, r)

		requestLogger.Info("request completed")
	})
}
