package tests_test

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	user2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

func TestAuthUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := user2.NewMockIUserRepository(ctrl)
	mockToken := user2.NewMockITokenator(ctrl)
	logger := logrus.New()
	minioConfig := &config.MinioConfig{
		Port:         "9000",
		Endpoint:     "localhost",
		BucketName:   "my-bucket",
		RootUser:     "minioadmin",
		RootPassword: "minioadminpassword",
		UseSSL:       false,
	}
	minio, err := minio.NewMinioClient(minioConfig)

	assert.Error(t, err)
	uc := user.NewAuthUsecase(mockRepo, mockToken, logger, minio)

	testUser := models.UserRegisterRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test",
		Surname:  null.StringFrom("User"),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			CheckUserExists(gomock.Any(), testUser.Email).
			Return(false, nil).
			Times(1)

		mockRepo.EXPECT().
			CreateUser(gomock.Any(), gomock.All(
				gomock.AssignableToTypeOf(models.UserDB{}),
				gomock.Not(gomock.Nil()),
			)).
			Return(nil).
			Times(1)

		mockToken.EXPECT().
			CreateJWT(gomock.Any(), gomock.Any()).
			Return("test-token", nil).
			Times(1)

		token, err := uc.Register(context.Background(), testUser)

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("CheckUserExistsError", func(t *testing.T) {
		mockRepo.EXPECT().
			CheckUserExists(gomock.Any(), testUser.Email).
			Return(false, errors.New("db error")).
			Times(1)

		token, err := uc.Register(context.Background(), testUser)

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("CreateUserError", func(t *testing.T) {
		mockRepo.EXPECT().
			CheckUserExists(gomock.Any(), testUser.Email).
			Return(false, nil).
			Times(1)

		mockRepo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(errors.New("db error")).
			Times(1)

		token, err := uc.Register(context.Background(), testUser)

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("CreateJWTError", func(t *testing.T) {
		mockRepo.EXPECT().
			CheckUserExists(gomock.Any(), testUser.Email).
			Return(false, nil).
			Times(1)

		mockRepo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		mockToken.EXPECT().
			CreateJWT(gomock.Any(), 1).
			Return("", errors.New("jwt error")).
			Times(1)

		token, err := uc.Register(context.Background(), testUser)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestAuthUsecase_Login(t *testing.T) {
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
	minio, err := minio.NewMinioClient(minioConfig)

	assert.Error(t, err)
	uc := user.NewAuthUsecase(mockRepo, mockToken, logger, minio)

	testUserID := uuid.New()

	testUser := models.UserLoginRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.MinCost)
	testUserDB := &models.UserDB{
		ID:           uuid.New(),
		Email:        testUser.Email,
		PasswordHash: hashedPassword,
	}

	t.Run("Success", func(t *testing.T) {
		userVersion := models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    testUserID,
			Version:   1,
			UpdatedAt: time.Now(),
		}

		userDB := &models.UserDB{
			ID:           testUserID,
			Email:        testUser.Email,
			PasswordHash: hashedPassword,
			UserVersion:  userVersion,
		}

		mockRepo.EXPECT().
			GetUserByEmail(gomock.Any(), testUser.Email).
			Return(userDB, nil).
			Times(1)

		mockToken.EXPECT().
			CreateJWT(userDB.ID.String(), userDB.UserVersion.Version).
			Return("test-token", nil).
			Times(1)

		token, err := uc.Login(context.Background(), testUser)

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUserByEmail(gomock.Any(), testUser.Email).
			Return(nil, errors.New("not found")).
			Times(1)

		token, err := uc.Login(context.Background(), testUser)

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		invalidUserDB := *testUserDB
		invalidUserDB.PasswordHash = []byte("wrong-hash")

		mockRepo.EXPECT().
			GetUserByEmail(gomock.Any(), testUser.Email).
			Return(&invalidUserDB, nil).
			Times(1)

		token, err := uc.Login(context.Background(), testUser)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidCredentials, err)
		assert.Empty(t, token)
	})

	t.Run("CreateJWTError", func(t *testing.T) {
		mockRepo.EXPECT().
			GetUserByEmail(gomock.Any(), testUser.Email).
			Return(testUserDB, nil).
			Times(1)

		mockToken.EXPECT().
			CreateJWT(testUserDB.ID.String(), testUserDB.UserVersion.Version).
			Return("", errors.New("jwt error")).
			Times(1)

		token, err := uc.Login(context.Background(), testUser)

		assert.Error(t, err)
		assert.Empty(t, token)
	})

}

func TestAuthUsecase_Logout(t *testing.T) {
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
	minio, err := minio.NewMinioClient(minioConfig)

	assert.Error(t, err)
	uc := user.NewAuthUsecase(mockRepo, mockToken, logger, minio)

	testUserID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), cookie.UserIDKey, testUserID)

		mockRepo.EXPECT().
			IncrementUserVersion(gomock.Any(), testUserID).
			Return(nil).
			Times(1)

		err := uc.Logout(ctx)

		assert.NoError(t, err)
	})

	t.Run("UserNotFoundInContext", func(t *testing.T) {
		err := uc.Logout(context.Background())

		assert.Error(t, err)
		assert.Equal(t, errs.ErrUserNotFound, err)
	})

	t.Run("IncrementVersionError", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), cookie.UserIDKey, testUserID)

		mockRepo.EXPECT().
			IncrementUserVersion(gomock.Any(), testUserID).
			Return(errors.New("db error")).
			Times(1)

		err := uc.Logout(ctx)

		assert.Error(t, err)
	})
}

func TestAuthUsecase_GetMe(t *testing.T) {
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
	minio, err := minio.NewMinioClient(minioConfig)

	assert.Error(t, err)
	uc := user.NewAuthUsecase(mockRepo, mockToken, logger, minio)

	testUserID := uuid.New()
	testUserDB := &models.UserDB{
		ID:      testUserID,
		Email:   "test@example.com",
		Name:    "Test",
		Surname: null.StringFrom("User"),
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), cookie.UserIDKey, testUserID.String())

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
		assert.Equal(t, errs.ErrUserNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), cookie.UserIDKey, "invalid-uuid")

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrInvalidUserID, err)
		assert.Nil(t, user)
	})

	t.Run("GetUserByIDError", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), cookie.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(nil, errors.New("db error")).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("UserNotFoundInDB", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), cookie.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(nil, nil).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Equal(t, errs.ErrUserNotFound, err)
		assert.Nil(t, user)
	})
}
