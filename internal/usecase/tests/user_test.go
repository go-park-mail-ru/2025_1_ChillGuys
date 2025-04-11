package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	user2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserUsecase_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
	testUserDB := &dto.UserDB{
		ID:      testUserID,
		Email:   "test@example.com",
		Name:    "Test",
		Surname: null.StringFrom("User"),
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(testUserDB, nil).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.NoError(t, err)
		assert.Equal(t, testUserDB.Email, user.Email)
		assert.Equal(t, testUserDB.Name, user.Name)
		assert.Equal(t, testUserDB.Surname, user.Surname)
	})

	t.Run("UserNotFoundInContext", func(t *testing.T) {
		user, err := uc.GetMe(context.Background())

		assert.Error(t, err)
		assert.Equal(t, errs.ErrNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey, "invalid-uuid")

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidID, err)
		assert.Nil(t, user)
	})

	t.Run("GetUserByIDError", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(nil, errors.New("db error")).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("UserNotFoundInDB", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domains.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(nil, nil).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrNotFound, err)
		assert.Nil(t, user)
	})
}
