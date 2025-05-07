package tests

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	authhttp "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth/http"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	genmock "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := genmock.NewMockAuthServiceClient(ctrl)

	// Создаем конфигурацию
	cfg := &config.Config{
		CSRFConfig: &config.CSRFConfig{SecretKey: "secret-key", TokenExpiry: time.Hour},
	}

	handler := authhttp.NewAuthHandler(mockClient, cfg)

	t.Run("success", func(t *testing.T) {
		reqData := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqData)

		// Мокаем ответ GRPC
		loginRes := &gen.LoginRes{Token: "jwt-token"}
		mockClient.EXPECT().Login(gomock.Any(), &gen.LoginReq{
			Email:    reqData["email"],
			Password: reqData["password"],
		}).Return(loginRes, nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		// Используем функцию Login из обработчика
		handler.Login(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Set-Cookie"), domains.TokenCookieName)
		assert.NotEmpty(t, resp.Header.Get("X-CSRF-Token"))
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("login failure", func(t *testing.T) {
		reqData := map[string]string{
			"email":    "bad@example.com",
			"password": "wrong",
		}
		body, _ := json.Marshal(reqData)

		mockClient.EXPECT().Login(gomock.Any(), gomock.Any()).
			Return(nil, status.Error(codes.Unauthenticated, "unauthorized"))

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := genmock.NewMockAuthServiceClient(ctrl)

	// Create config
	cfg := &config.Config{
		CSRFConfig: &config.CSRFConfig{SecretKey: "secret-key", TokenExpiry: time.Hour},
	}

	handler := authhttp.NewAuthHandler(mockClient, cfg)

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("user already exists", func(t *testing.T) {
		reqData := map[string]string{
			"email":    "existinguser@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqData)

		mockClient.EXPECT().Register(gomock.Any(), gomock.Any()).
			Return(nil, status.Error(codes.AlreadyExists, "user already exists"))

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := genmock.NewMockAuthServiceClient(ctrl)

	// Create config
	cfg := &config.Config{
		CSRFConfig: &config.CSRFConfig{SecretKey: "secret-key", TokenExpiry: time.Hour},
	}

	handler := authhttp.NewAuthHandler(mockClient, cfg)

	t.Run("success", func(t *testing.T) {
		// Mock JWT token
		jwtToken := "jwt-token"
		mockClient.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		req.AddCookie(&http.Cookie{Name: string(domains.TokenCookieName), Value: jwtToken})
		w := httptest.NewRecorder()

		// Use Logout handler
		handler.Logout(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Set-Cookie"), domains.TokenCookieName)
	})

	t.Run("no jwt token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("logout failure", func(t *testing.T) {
		jwtToken := "jwt-token"
		mockClient.EXPECT().Logout(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unauthenticated, "unauthorized"))

		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		req.AddCookie(&http.Cookie{Name: string(domains.TokenCookieName), Value: jwtToken})
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})
}
