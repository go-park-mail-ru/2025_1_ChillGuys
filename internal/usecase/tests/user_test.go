package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	user2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserUsecase_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	
	// Since Minio is not critical for auth tests, we can pass nil
	uc := user.NewAuthUsecase(mockRepo, mockToken, nil)
	
	return mockRepo, mockToken, uc
}

	mockRepo := user2.NewMockIUserRepository(ctrl)
	mockToken := user2.NewMockITokenator(ctrl)
	logger := logrus.New()
	minioConfig := &config.MinioConfig{
		Port:         "9000",               // Порт Minio
		Endpoint:     "localhost",          // Адрес Minio
		BucketName:   "my-bucket",          // Название бакета
		RootUser:     "minioadmin",         // Имя пользователя
		RootPassword: "minioadminpassword", // Пароль пользователя
		UseSSL:       false,                // Не используем SSL
	}

	ctx := context.Background()
	minio, err := minio.NewMinioClient(ctx, minioConfig)

	assert.Error(t, err)
	uc := user.NewUserUsecase(mockRepo, mockToken, logger, minio)

	testUserID := uuid.New()
	testUserDB := &models.UserDB{
		ID:      testUserID,
		Email:   "test@example.com",
		Name:    "Test",
		Surname: null.StringFrom("User"),
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey{}, testUserID.String())

		mockRepo.EXPECT().
			GetByID(gomock.Any(), userID).
			Return(userDB, nil)

		user, err := uc.GetMe(ctx)
		assert.NoError(t, err)
		assert.Equal(t, userDB.ConvertToUser(), user)
	})

	t.Run("no user in context", func(t *testing.T) {
		_, _, uc := setupTest(t)

		_, err := uc.GetMe(context.Background())
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey{}, "invalid-uuid")

		_, err := uc.GetMe(ctx)
		assert.ErrorIs(t, err, errs.ErrInvalidID)
	})

	t.Run("GetUserByIDError", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey{}, testUserID.String())

		mockRepo.EXPECT().
			GetByID(gomock.Any(), userID).
			Return(nil, errs.ErrNotFound)

		_, err := uc.GetMe(ctx)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("UserNotFoundInDB", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey{}, testUserID.String())

		mockRepo.EXPECT().
			GetByID(gomock.Any(), userID).
			Return(nil, errors.New("db error"))

		_, err := uc.GetMe(ctx)
		assert.Error(t, err)
	})
}