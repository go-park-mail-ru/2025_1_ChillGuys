package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t) // Создание контроллера для моков
	defer ctrl.Finish()             // Завершаем контроллер после выполнения теста

	// Создаем моки для репозитория пользователей и токенизатора
	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockTokenator := mocks.NewMockITokenator(ctrl)
	logger := logrus.New()

	// Создаем обработчик с моками
	handler := transport.NewAuthHandler(mockRepo, logger, mockTokenator)

	// Генерация хеша пароля для тестов
	passwordHash, err := transport.GeneratePasswordHash("Password123")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":    err,
			"password": "Password123",
		}).Error("Failed to generate password hash")
	}

	tests := []struct {
		name           string
		request        models.UserLoginRequestDTO
		mockBehavior   func()
		expectedStatus int
		expectedBody   null.String
	}{
		// Тест для успешного входа
		{
			name: "Valid Login",
			request: models.UserLoginRequestDTO{
				Email:    "test@example.com",
				Password: "Password123",
			},
			mockBehavior: func() {
				// Мокируем успешный поиск пользователя и создание JWT токена
				mockRepo.EXPECT().GetUserByEmail("test@example.com").Return(&models.UserDB{
					ID:           uuid.New(),
					Email:        "test@example.com",
					PasswordHash: passwordHash,
					Version:      1,
				}, nil)

				mockTokenator.EXPECT().
					CreateJWT(gomock.Any(), gomock.Any()).
					Return("mocked-jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   null.String{},
		},
		// Тест для некорректного формата email
		{
			name: "Invalid Email Format",
			request: models.UserLoginRequestDTO{
				Email:    "invalid-email",
				Password: "Password123",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"Invalid email"}`),
		},
		// Тест для некорректного пароля
		{
			name: "Invalid Password",
			request: models.UserLoginRequestDTO{
				Email:    "test@example.com",
				Password: "short",
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"Invalid password: password must be at least 8 characters"}`),
		},
		// Тест для случая, когда пользователь не найден
		{
			name: "User Not Found",
			request: models.UserLoginRequestDTO{
				Email:    "notfound@example.com",
				Password: "Password123",
			},
			mockBehavior: func() {
				mockRepo.EXPECT().GetUserByEmail("notfound@example.com").Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   null.StringFrom(`{"message":"Invalid password or email"}`),
		},
	}

	// Запускаем тесты
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior() // Выполняем мокируемое поведение

			// Создаем запрос с телом запроса в формате JSON
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder() // Записываем в рекордер ответ

			handler.Login(w, req)

			// Проверяем статус и тело ответа
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody.Valid {
				assert.JSONEq(t, tt.expectedBody.String, w.Body.String())
			} else {
				// Достаем токен из cookies
				var token string
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						token = cookie.Value
						break
					}
				}
				// Проверяем, что токен не пустой
				assert.NotEmpty(t, token, "Token should not be empty")
				assert.Equal(t, "mocked-jwt-token", token, "JWT token does not match expected")
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t) // Создание контроллера для моков
	defer ctrl.Finish()             // Завершаем контроллер после выполнения теста

	// Создаем моки для репозитория пользователей и токенизатора
	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockTokenator := mocks.NewMockITokenator(ctrl)
	logger := logrus.New()

	// Создаем обработчик с моками
	handler := transport.NewAuthHandler(mockRepo, logger, mockTokenator)

	// Тесты для различных сценариев
	tests := []struct {
		name           string
		request        models.UserRegisterRequestDTO // Структура запроса
		mockBehavior   func()                        // Мокируемое поведение
		expectedStatus int                           // Ожидаемый статус ответа
		expectedBody   null.String                   // Ожидаемое тело ответа
	}{
		// Тест для успешной регистрации
		{
			name: "Valid Registration",
			request: models.UserRegisterRequestDTO{
				Email:    "newuser@example.com",
				Password: "Password123",
				Name:     "John",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior: func() {
				// Мокируем успешную проверку, что пользователь не существует, и создание нового пользователя
				mockRepo.EXPECT().GetUserByEmail("newuser@example.com").Return(nil, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
				mockTokenator.EXPECT().CreateJWT(gomock.Any(), gomock.Any()).Return("mocked-jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   null.String{},
		},
		// Тест для некорректного email
		{
			name: "Invalid Email",
			request: models.UserRegisterRequestDTO{
				Email:    "invalid-email",
				Password: "Password123",
				Name:     "John",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"Invalid email"}`),
		},
		// Тест для уже существующего пользователя
		{
			name: "User Already Exists",
			request: models.UserRegisterRequestDTO{
				Email:    "existing@example.com",
				Password: "Password123",
				Name:     "John",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior: func() {
				// Мокируем, что пользователь с таким email уже существует
				mockRepo.EXPECT().GetUserByEmail("existing@example.com").Return(&models.UserDB{}, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   null.StringFrom(`{"message":"User already exists"}`),
		},
		// Тест для пустого имени
		{
			name: "Empty Name",
			request: models.UserRegisterRequestDTO{
				Email:    "newuser@example.com",
				Password: "Password123",
				Name:     "",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"Invalid name"}`),
		},
		// Тест для некорректного пароля
		{
			name: "Short Password",
			request: models.UserRegisterRequestDTO{
				Email:    "newuser@example.com",
				Password: "Pass",
				Name:     "John",
				Surname:  null.StringFrom("Doe"),
			},
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   null.StringFrom(`{"message":"Invalid password: password must be at least 8 characters"}`),
		},
	}

	// Запускаем тесты
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior() // Выполняем мокируемое поведение

			// Создаем запрос с телом запроса в формате JSON
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder() // Записываем в рекордер ответ

			handler.Register(w, req)

			// Проверяем статус и тело ответа
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody.Valid {
				assert.JSONEq(t, tt.expectedBody.String, w.Body.String())
			} else {
				// Достаем токен из cookies
				var token string
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "token" {
						token = cookie.Value
						break
					}
				}
				// Проверяем, что токен не пустой
				assert.NotEmpty(t, token, "Token should not be empty")
				assert.Equal(t, "mocked-jwt-token", token, "JWT token does not match expected")
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t) // Создание контроллера для моков
	defer ctrl.Finish()             // Завершаем контроллер после выполнения теста

	// Создаем моки для репозитория пользователей, токенизатора
	mockRepo := mocks.NewMockIUserRepository(ctrl)
	logger := logrus.New()
	mockTokenator := mocks.NewMockITokenator(ctrl)
	handler := transport.NewAuthHandler(mockRepo, logger, mockTokenator)

	// Тесты для различных сценариев
	tests := []struct {
		name           string
		userID         string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		// Тест для успешного выхода
		{
			name:   "Successful Logout",
			userID: "user-id",
			mockBehavior: func() {
				mockRepo.EXPECT().IncrementUserVersion("user-id").Return(nil) // Мокируем успешное обновление версии пользователя
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`, // Ожидаем пустое тело в ответе
		},
		// Тест для случая, когда не передан ID пользователя
		{
			name:           "User ID Not Found",
			userID:         "",
			mockBehavior:   func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"User id not found"}`,
		},
		// Тест для ошибки при обновлении версии пользователя
		{
			name:   "Increment Version Failed",
			userID: "user-id",
			mockBehavior: func() {
				mockRepo.EXPECT().IncrementUserVersion("user-id").Return(errors.New("database error")) // Мокируем ошибку в базе данных
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database error"}`,
		},
	}

	// Запускаем тесты
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior() // Выполняем мокируемое поведение

			// Создаем запрос для выхода
			req := httptest.NewRequest("POST", "/logout", nil)
			if tt.userID != "" {
				req = req.WithContext(context.WithValue(req.Context(), "userID", tt.userID)) // Устанавливаем userID в контекст
			}
			w := httptest.NewRecorder() // Записываем в рекордер ответ

			handler.Logout(w, req)

			// Проверяем статус и тело ответа
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody == `{}` {
				assert.Empty(t, w.Body.String()) // Проверяем, что тело пустое для успешного выхода
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String()) // Проверяем тело ответа на совпадение с ожидаемым
			}
		})
	}
}

func TestUserHandler_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t) // Создание контроллера для моков
	defer ctrl.Finish()             // Завершаем контроллер после выполнения теста

	// Создаем моки для репозитория пользователей и токенизатора
	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockTokenator := mocks.NewMockITokenator(ctrl)
	logger := logrus.New()

	// Создаем обработчик с моками
	handler := transport.NewAuthHandler(mockRepo, logger, mockTokenator)

	// Генерация тестовых данных
	userID := uuid.New()
	userVersion := 1

	// Тесты для различных сценариев
	tests := []struct {
		name           string
		userID         string
		userVersion    int
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		// Тест для успешного получения информации о пользователе
		{
			name:        "Successful GetMe",
			userID:      userID.String(),
			userVersion: userVersion,
			mockBehavior: func() {
				// Мокируем успешное получение пользователя
				mockRepo.EXPECT().GetUserByID(userID).Return(&models.UserDB{
					ID:          userID,
					Email:       "test@example.com",
					Name:        "John",
					Surname:     null.StringFrom("Doe"),
					PhoneNumber: null.StringFrom("1234567890"),
					Version:     userVersion,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"` + userID.String() + `","email":"test@example.com","name":"John","surname":"Doe","phoneNumber":"1234567890"}`,
		},
		// Тест для случая, когда не передан ID пользователя в контексте
		{
			name:           "User ID Not Found",
			userID:         "",
			userVersion:    userVersion,
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid user id format"}`,
		},
		// Тест для случая с ошибкой при получении пользователя
		{
			name:           "User ID Not Found",
			userID:         "",
			mockBehavior:   func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid user id format"}`,
		},
	}

	// Запускаем тесты
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior() // Выполняем мокируемое поведение

			// Создаем запрос для получения информации о пользователе
			req := httptest.NewRequest("GET", "/me", nil)
			req = req.WithContext(context.WithValue(req.Context(), "userID", tt.userID))
			req = req.WithContext(context.WithValue(req.Context(), "userVersion", tt.userVersion))

			w := httptest.NewRecorder() // Записываем в рекордер ответ

			handler.GetMe(w, req)

			// Проверяем статус и тело ответа
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody == `{}` {
				assert.Empty(t, w.Body.String()) // Проверяем, что тело пустое для успешного выхода
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String()) // Проверяем тело ответа на совпадение с ожидаемым
			}
		})
	}
}
