package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testHook struct {
	entries []*logrus.Entry
}

func (h *testHook) Fire(entry *logrus.Entry) error {
	h.entries = append(h.entries, entry)
	return nil
}

func (h *testHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func TestLogRequest(t *testing.T) {
	// Создаем тестовый логгер с хуком
	logger := logrus.New()
	hook := &testHook{}
	logger.AddHook(hook)

	// Создаем тестовый обработчик, который проверит контекст
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие request_id в контексте
		reqID, ok := r.Context().Value(domains.ReqIDKey{}).(string)
		require.True(t, ok, "Request ID should be in context")
		assert.NotEmpty(t, reqID, "Request ID should not be empty")

		// Проверяем наличие логгера в контексте
		logger := logctx.GetLogger(r.Context())
		assert.NotNil(t, logger, "Logger should be in context")
		assert.Equal(t, reqID, logger.Data["request_id"], "Logger should have correct request_id")

		w.WriteHeader(http.StatusOK)
	})

	// Создаем middleware с тестовым логгером
	middleware := middleware.LogRequest(logger, nextHandler)

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	rec := httptest.NewRecorder()

	// Выполняем запрос
	startTime := time.Now()
	middleware.ServeHTTP(rec, req)
	duration := time.Since(startTime)

	// Проверяем результаты
	assert.Equal(t, http.StatusOK, rec.Code, "Handler should return 200 OK")

	// Проверяем, что было как минимум 2 записи в лог (начало и конец)
	require.GreaterOrEqual(t, len(hook.entries), 2, "Should have at least 2 log entries")

	// Проверяем первую запись (начало обработки)
	firstEntry := hook.entries[0]
	assert.Equal(t, "request started", firstEntry.Message)
	assert.Equal(t, "GET", firstEntry.Data["method"])
	assert.Equal(t, "127.0.0.1:12345", firstEntry.Data["remote_addr"])
	assert.Equal(t, "/test", firstEntry.Data["path"])
	assert.NotEmpty(t, firstEntry.Data["request_id"])

	// Проверяем последнюю запись (конец обработки)
	lastEntry := hook.entries[len(hook.entries)-1]
	assert.Equal(t, "request completed", lastEntry.Message)
	assert.NotEmpty(t, lastEntry.Data["duration"])
	assert.Equal(t, firstEntry.Data["request_id"], lastEntry.Data["request_id"])

	// Проверяем duration
	logDuration, ok := lastEntry.Data["duration"].(time.Duration)
	require.True(t, ok, "Duration should be time.Duration")
	assert.True(t, logDuration > 0, "Duration should be positive")
	assert.True(t, logDuration <= duration, "Logged duration should be less than or equal to actual duration")
}

func TestLogRequest_ContextValues(t *testing.T) {
	logger := logrus.New()
	hook := &testHook{}
	logger.AddHook(hook)

	var ctx context.Context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		w.WriteHeader(http.StatusOK)
	})

	middleware := middleware.LogRequest(logger, nextHandler)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Проверяем значения в контексте
	reqID, ok := ctx.Value(domains.ReqIDKey{}).(string)
	require.True(t, ok, "Request ID should be in context")
	assert.NotEmpty(t, reqID, "Request ID should not be empty")

	loggerEntry, ok := ctx.Value(domains.LoggerKey{}).(*logrus.Entry)
	require.True(t, ok, "Logger should be in context")
	assert.Equal(t, reqID, loggerEntry.Data["request_id"], "Logger should have correct request_id")

	// Проверяем что request_id в логгере совпадает с request_id в контексте
	require.GreaterOrEqual(t, len(hook.entries), 1, "Should have log entries")
	assert.Equal(t, reqID, hook.entries[0].Data["request_id"], "Logger entry should match context request_id")
}

// func TestLogRequest_NoRemoteAddr(t *testing.T) {
// 	logger := logrus.New()
// 	hook := &testHook{}
// 	logger.AddHook(hook)

// 	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	})

// 	middleware := middleware.LogRequest(logger, nextHandler)

// 	req := httptest.NewRequest("GET", "http://example.com/test", nil)
// 	req.RemoteAddr = "" // Пустой RemoteAddr

// 	rec := httptest.NewRecorder()

// 	middleware.ServeHTTP(rec, req)

// 	require.GreaterOrEqual(t, len(hook.entries), 1, "Should have log entries")
// 	assert.NotEmpty(t, hook.entries[0].Data["remote_addr"], "Should handle empty remote address")
// }

func TestLogRequest_LoggingOrder(t *testing.T) {
	logger := logrus.New()
	hook := &testHook{}
	logger.AddHook(hook)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Имитируем обработку
		w.WriteHeader(http.StatusOK)
	})

	middleware := middleware.LogRequest(logger, nextHandler)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Проверяем что было как минимум 2 записи (начало и конец)
	require.GreaterOrEqual(t, len(hook.entries), 2, "Should have start and end log entries")
	assert.Equal(t, "request started", hook.entries[0].Message)
	assert.Equal(t, "request completed", hook.entries[len(hook.entries)-1].Message)
}
