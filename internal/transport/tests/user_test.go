package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logrus.New()
	mockUserUsecase := mocks.NewMockIUserUsecase(ctrl)

	minioConfig := &config.MinioConfig{
		Port:         "9000",
		Endpoint:     "localhost",
		BucketName:   "my-bucket",
		RootUser:     "minioadmin",
		RootPassword: "minioadminpassword",
		UseSSL:       false,
	}

	ctx := context.Background()
	minio, err := minio.NewMinioClient(ctx, minioConfig)
	assert.Error(t, err)

	handler := user.NewUserHandler(mockUserUsecase, logger, minio, &config.Config{})

	userID := uuid.New()

	tests := []struct {
		name           string
		userID         string
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful GetMe",
			userID: userID.String(),
			mockBehavior: func() {
				mockUserUsecase.EXPECT().
					GetMe(gomock.Any()).
					Return(&dto.UserDTO{
						ID:          userID,
						Email:       "test@example.com",
						Name:        "John",
						Surname:     null.StringFrom("Doe"),
						PhoneNumber: null.StringFrom("1234567890"),
						ImageURL:    null.String{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"` + userID.String() + `","email":"test@example.com","name":"John","surname":"Doe","phoneNumber":"1234567890","imageURL":null}`,
		},
		{
			name:   "User Not Found",
			userID: userID.String(),
			mockBehavior: func() {
				mockUserUsecase.EXPECT().
					GetMe(gomock.Any()).
					Return((*dto.UserDTO)(nil), errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"user not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := httptest.NewRequest("GET", "/me", nil)
			if tt.userID != "" {
				req = req.WithContext(context.WithValue(req.Context(), "userID", tt.userID))
			}
			w := httptest.NewRecorder()

			handler.GetMe(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody == `{}` {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
