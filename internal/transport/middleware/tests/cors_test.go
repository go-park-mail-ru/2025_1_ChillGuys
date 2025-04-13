package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
)

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	// Создаем тестовую конфигурацию
	conf := &config.ServerConfig{
		AllowOrigin:      "http://allowed.com",
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Content-Type",
		AllowCredentials: "true",
	}

	// Создаем тестовый обработчик
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := middleware.CORSMiddleware(nextHandler, conf)

	tests := []struct {
		name           string
		origin         string
		expectedStatus int
		expectedOrigin string
	}{
		{
			name:           "Allowed origin",
			origin:         "http://allowed.com",
			expectedStatus: http.StatusOK,
			expectedOrigin: "http://allowed.com",
		},
		{
			name:           "Not allowed origin",
			origin:         "http://not-allowed.com",
			expectedStatus: http.StatusForbidden,
			expectedOrigin: "",
		},
		{
			name:           "No origin header",
			origin:         "",
			expectedStatus: http.StatusOK,
			expectedOrigin: "", // Для не-CORS запросов заголовок не должен устанавливаться
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			originHeader := rec.Header().Get("Access-Control-Allow-Origin")
			if tt.expectedOrigin != "" && originHeader != tt.expectedOrigin {
				t.Errorf("expected origin header '%s', got '%s'", tt.expectedOrigin, originHeader)
			}
			if tt.expectedOrigin == "" && originHeader != "" {
				t.Errorf("expected no origin header, got '%s'", originHeader)
			}
		})
	}
}