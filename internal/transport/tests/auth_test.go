package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	http2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth/http"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockIAuthUsecase(ctrl)
	handler := http2.NewAuthHandler(mockAuthUsecase, &config.Config{})

	tests := []struct {
		name           string
		request        dto.UserLoginRequestDTO
		mockBehavior   func()
		expectedStatus int
		expectedBody   null.String
	}{
		//{
		//	name: "Valid Login",
		//	request: dto.UserLoginRequestDTO{
		//		Email:    "test@example.com",
		//		Password: "Password123",
		//	},
		//	mockBehavior: func() {
		//		mockAuthUsecase.EXPECT().
		//			Login(gomock.Any(), gomock.Any()).
		//			Return("mocked-jwt-token", nil)
		//	},
		//	expectedStatus: http.StatusOK,
		//	expectedBody:   null.String{},
		//},
		{
			name: "Invalid Email Format",
			request: dto.UserLoginRequestDTO{
				Email:    "invalid-email",
				Password: "Password123",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"invalid email"}`),
		},
		{
			name: "Invalid Password",
			request: dto.UserLoginRequestDTO{
				Email:    "test@example.com",
				Password: "short",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"password must be at least 8 characters"}`),
		},
		//{
		//	name: "User Not Found",
		//	request: dto.UserLoginRequestDTO{
		//		Email:    "notfound@example.com",
		//		Password: "Password123",
		//	},
		//	mockBehavior: func() {
		//		mockAuthUsecase.EXPECT().
		//			Login(gomock.Any(), gomock.Any()).
		//			Return("", errors.New("auth not found"))
		//	},
		//	expectedStatus: http.StatusInternalServerError,
		//	expectedBody:   null.StringFrom(`{"message":"auth not found"}`),
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Упрощенная проверка ответа
			if !tt.expectedBody.Valid || tt.expectedBody.String == "" {
				// Для успешного логина проверяем только статус
				assert.Equal(t, http.StatusOK, w.Code)
				// Можно добавить проверку, что тело пустое
				assert.Empty(t, w.Body.String())

				// Если нужно проверить куки (опционально):
				cookies := w.Result().Cookies()
				assert.NotEmpty(t, cookies, "Expected cookies to be set")
				// Дополнительные проверки кук если нужно
			} else {
				// Для ошибок проверяем JSON тело
				assert.JSONEq(t, tt.expectedBody.String, w.Body.String())
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockIAuthUsecase(ctrl)

	handler := http2.NewAuthHandler(mockAuthUsecase, &config.Config{})

	tests := []struct {
		name           string
		request        dto.UserRegisterRequestDTO
		mockBehavior   func()
		expectedStatus int
		expectedBody   null.String
	}{
		//{
		//	name: "Valid Registration",
		//	request: dto.UserRegisterRequestDTO{
		//		Email:    "newuser@example.com",
		//		Password: "Password123",
		//		Name:     "John",
		//		Surname:  null.StringFrom("Doe"),
		//	},
		//	mockBehavior: func() {
		//		mockAuthUsecase.EXPECT().
		//			Register(gomock.Any(), gomock.Any()).
		//			Return("mocked-jwt-token", nil)
		//	},
		//	expectedStatus: http.StatusOK,
		//	expectedBody:   null.String{},
		//},
		{
			name: "Invalid Email",
			request: dto.UserRegisterRequestDTO{
				Email:    "invalid-email",
				Password: "Password123",
				Name:     "John",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"invalid email"}`),
		},
		//{
		//	name: "User Already Exists",
		//	request: dto.UserRegisterRequestDTO{
		//		Email:    "existing@example.com",
		//		Password: "Password123",
		//		Name:     "John",
		//		Surname:  null.StringFrom("Doe"),
		//	},
		//	mockBehavior: func() {
		//		mockAuthUsecase.EXPECT().
		//			Register(gomock.Any(), gomock.Any()).
		//			Return("", errors.New("auth already exists"))
		//	},
		//	expectedStatus: http.StatusInternalServerError,
		//	expectedBody:   null.StringFrom(`{"message":"auth already exists"}`),
		//},
		{
			name: "Empty Name",
			request: dto.UserRegisterRequestDTO{
				Email:    "newuser@example.com",
				Password: "Password123",
				Name:     "",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"name must be between 2 and 24 characters long"}`),
		},
		{
			name: "Short Password",
			request: dto.UserRegisterRequestDTO{
				Email:    "newuser@example.com",
				Password: "Pass",
				Name:     "John",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"password must be at least 8 characters"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Register(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Правильная проверка с использованием null.String
			if !tt.expectedBody.Valid {
				// Для случая когда expectedBody не валиден (ожидаем успешный ответ без тела)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Empty(t, w.Body.String())
			} else {
				// Для случаев с ошибками (ожидаем JSON в теле ответа)
				assert.JSONEq(t, tt.expectedBody.String, w.Body.String())
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := mocks.NewMockIAuthUsecase(ctrl)

	handler := http2.NewAuthHandler(mockAuthUsecase, &config.Config{})

	tests := []struct {
		name           string
		userID         string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful Logout",
			userID: "auth-id",
			mockBehavior: func() {
				mockAuthUsecase.EXPECT().
					Logout(gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`,
		},
		{
			name:   "Increment Version Failed",
			userID: "auth-id",
			mockBehavior: func() {
				mockAuthUsecase.EXPECT().
					Logout(gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("POST", "/logout", nil)
			if tt.userID != "" {
				req = req.WithContext(context.WithValue(req.Context(), "userID", tt.userID))
			}
			w := httptest.NewRecorder()

			handler.Logout(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody == `{}` {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
