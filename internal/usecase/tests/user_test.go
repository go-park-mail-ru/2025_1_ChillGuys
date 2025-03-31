package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase"
)

func TestAuthUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	logger := logrus.New()
	uc := usecase.NewAuthUsecase(mockRepo, mockToken, logger)

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
			CreateUser(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, user models.UserDB) error {
				assert.Equal(t, testUser.Email, user.Email)
				assert.Equal(t, testUser.Name, user.Name)
				assert.Equal(t, testUser.Surname, user.Surname)
				assert.NotEmpty(t, user.PasswordHash)
				assert.Equal(t, 1, user.Version)
				return nil
			}).
			Times(1)

		mockToken.EXPECT().
			CreateJWT(gomock.Any(), 1).
			Return("test-token", nil).
			Times(1)

		token, err := uc.Register(context.Background(), testUser)

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mockRepo.EXPECT().
			CheckUserExists(gomock.Any(), testUser.Email).
			Return(true, nil).
			Times(1)

		token, err := uc.Register(context.Background(), testUser)

		assert.Error(t, err)
		assert.Equal(t, models.ErrUserAlreadyExists, err)
		assert.Empty(t, token)
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

    mockRepo := mocks.NewMockIUserRepository(ctrl)
    mockToken := mocks.NewMockITokenator(ctrl)
    logger := logrus.New()
    uc := usecase.NewAuthUsecase(mockRepo, mockToken, logger)

    testUser := models.UserLoginRequestDTO{
        Email:    "test@example.com",
        Password: "password123",
    }

    // Генерируем реальный хэш для тестов
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.MinCost)
    testUserDB := &models.UserDB{
        ID:           uuid.New(),
        Email:        testUser.Email,
        PasswordHash: hashedPassword,
        Version:      1,
    }

    t.Run("Success", func(t *testing.T) {
        mockRepo.EXPECT().
            GetUserByEmail(gomock.Any(), testUser.Email).
            Return(testUserDB, nil).
            Times(1)

        mockToken.EXPECT().
            CreateJWT(testUserDB.ID.String(), testUserDB.Version).
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
        // Создаем пользователя с неверным паролем
        invalidUserDB := *testUserDB
        invalidUserDB.PasswordHash = []byte("wrong-hash")

        mockRepo.EXPECT().
            GetUserByEmail(gomock.Any(), testUser.Email).
            Return(&invalidUserDB, nil).
            Times(1)

        token, err := uc.Login(context.Background(), testUser)

        assert.Error(t, err)
        assert.Equal(t, models.ErrInvalidCredentials, err)
        assert.Empty(t, token)
    })

    t.Run("CreateJWTError", func(t *testing.T) {
        mockRepo.EXPECT().
            GetUserByEmail(gomock.Any(), testUser.Email).
            Return(testUserDB, nil).
            Times(1)

        mockToken.EXPECT().
            CreateJWT(testUserDB.ID.String(), testUserDB.Version).
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

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	logger := logrus.New()
	uc := usecase.NewAuthUsecase(mockRepo, mockToken, logger)

	testUserID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, testUserID)

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
		assert.Equal(t, models.ErrUserNotFound, err)
	})

	t.Run("IncrementVersionError", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, testUserID)

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

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	logger := logrus.New()
	uc := usecase.NewAuthUsecase(mockRepo, mockToken, logger)

	testUserID := uuid.New()
	testUserDB := &models.UserDB{
		ID:      testUserID,
		Email:   "test@example.com",
		Name:    "Test",
		Surname: null.StringFrom("User"),
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, testUserID.String())

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
		assert.Equal(t, models.ErrUserNotFound, err)
		assert.Nil(t, user)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, "invalid-uuid")

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Equal(t, models.ErrInvalidUserID, err)
		assert.Nil(t, user)
	})

	t.Run("GetUserByIDError", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(nil, errors.New("db error")).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("UserNotFoundInDB", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), utils.UserIDKey, testUserID.String())

		mockRepo.EXPECT().
			GetUserByID(gomock.Any(), testUserID).
			Return(nil, nil).
			Times(1)

		user, err := uc.GetMe(ctx)

		assert.Error(t, err)
		assert.Equal(t, models.ErrUserNotFound, err)
		assert.Nil(t, user)
	})
}
