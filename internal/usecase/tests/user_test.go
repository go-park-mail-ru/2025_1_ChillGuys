package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupTest(t *testing.T) (*mocks.MockIUserRepository, *mocks.MockITokenator, *user.AuthUsecase) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	
	// Since Minio is not critical for auth tests, we can pass nil
	uc := user.NewAuthUsecase(mockRepo, mockToken, nil)
	
	return mockRepo, mockToken, uc
}

func TestAuthUsecase_Register(t *testing.T) {
	registerReq := dto.UserRegisterRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test",
		Surname:  null.StringFrom("User"),
	}

	t.Run("success", func(t *testing.T) {
		mockRepo, mockToken, uc := setupTest(t)

		mockRepo.EXPECT().
			CheckExistence(gomock.Any(), registerReq.Email).
			Return(false, nil)

		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, user dto.UserDB) error {
				assert.Equal(t, registerReq.Email, user.Email)
				assert.Equal(t, registerReq.Name, user.Name)
				assert.Equal(t, registerReq.Surname, user.Surname)
				assert.NotEmpty(t, user.PasswordHash)
				return nil
			})

		mockToken.EXPECT().
			CreateJWT(gomock.Any(), 1).
			Return("test-token", nil)

		token, err := uc.Register(context.Background(), registerReq)
		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)

		mockRepo.EXPECT().
			CheckExistence(gomock.Any(), registerReq.Email).
			Return(true, nil)

		_, err := uc.Register(context.Background(), registerReq)
		assert.ErrorIs(t, err, errs.ErrAlreadyExists)
	})

	t.Run("repository error on check existence", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)

		mockRepo.EXPECT().
			CheckExistence(gomock.Any(), registerReq.Email).
			Return(false, errors.New("db error"))

		_, err := uc.Register(context.Background(), registerReq)
		assert.Error(t, err)
	})

	t.Run("repository error on create", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)

		mockRepo.EXPECT().
			CheckExistence(gomock.Any(), registerReq.Email).
			Return(false, nil)

		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		_, err := uc.Register(context.Background(), registerReq)
		assert.Error(t, err)
	})

	t.Run("token creation error", func(t *testing.T) {
		mockRepo, mockToken, uc := setupTest(t)

		mockRepo.EXPECT().
			CheckExistence(gomock.Any(), registerReq.Email).
			Return(false, nil)

		mockRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil)

		mockToken.EXPECT().
			CreateJWT(gomock.Any(), 1).
			Return("", errors.New("token error"))

		_, err := uc.Register(context.Background(), registerReq)
		assert.Error(t, err)
	})
}

func TestAuthUsecase_Login(t *testing.T) {
	loginReq := dto.UserLoginRequestDTO{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginReq.Password), bcrypt.MinCost)
	userDB := &dto.UserDB{
		ID:           uuid.New(),
		Email:        loginReq.Email,
		PasswordHash: hashedPassword,
		UserVersion: models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo, mockToken, uc := setupTest(t)

		mockRepo.EXPECT().
			GetByEmail(gomock.Any(), loginReq.Email).
			Return(userDB, nil)

		mockToken.EXPECT().
			CreateJWT(userDB.ID.String(), userDB.UserVersion.Version).
			Return("test-token", nil)

		token, err := uc.Login(context.Background(), loginReq)
		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)

		mockRepo.EXPECT().
			GetByEmail(gomock.Any(), loginReq.Email).
			Return(nil, errs.ErrNotFound)

		_, err := uc.Login(context.Background(), loginReq)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)

		invalidUser := *userDB
		invalidUser.PasswordHash = []byte("wrong-hash")

		mockRepo.EXPECT().
			GetByEmail(gomock.Any(), loginReq.Email).
			Return(&invalidUser, nil)

		_, err := uc.Login(context.Background(), loginReq)
		assert.ErrorIs(t, err, errs.ErrInvalidCredentials)
	})

	t.Run("token creation error", func(t *testing.T) {
		mockRepo, mockToken, uc := setupTest(t)

		mockRepo.EXPECT().
			GetByEmail(gomock.Any(), loginReq.Email).
			Return(userDB, nil)

		mockToken.EXPECT().
			CreateJWT(userDB.ID.String(), userDB.UserVersion.Version).
			Return("", errors.New("token error"))

		_, err := uc.Login(context.Background(), loginReq)
		assert.Error(t, err)
	})
}

func TestAuthUsecase_Logout(t *testing.T) {
	userID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)
		ctx := context.WithValue(context.Background(), domains.UserIDKey, userID)

		mockRepo.EXPECT().
			IncrementVersion(gomock.Any(), userID).
			Return(nil)

		err := uc.Logout(ctx)
		assert.NoError(t, err)
	})

	t.Run("no user in context", func(t *testing.T) {
		_, _, uc := setupTest(t)

		err := uc.Logout(context.Background())
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)
		ctx := context.WithValue(context.Background(), domains.UserIDKey, userID)

		mockRepo.EXPECT().
			IncrementVersion(gomock.Any(), userID).
			Return(errors.New("db error"))

		err := uc.Logout(ctx)
		assert.Error(t, err)
	})
}

func TestAuthUsecase_GetMe(t *testing.T) {
	userID := uuid.New()
	userDB := &dto.UserDB{
		ID:      userID,
		Email:   "test@example.com",
		Name:    "Test",
		Surname: null.StringFrom("User"),
	}

	t.Run("success", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)
		ctx := context.WithValue(context.Background(), domains.UserIDKey, userID.String())

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

	t.Run("invalid user id format", func(t *testing.T) {
		_, _, uc := setupTest(t)
		ctx := context.WithValue(context.Background(), domains.UserIDKey, "invalid-uuid")

		_, err := uc.GetMe(ctx)
		assert.ErrorIs(t, err, errs.ErrInvalidID)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)
		ctx := context.WithValue(context.Background(), domains.UserIDKey, userID.String())

		mockRepo.EXPECT().
			GetByID(gomock.Any(), userID).
			Return(nil, errs.ErrNotFound)

		_, err := uc.GetMe(ctx)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, _, uc := setupTest(t)
		ctx := context.WithValue(context.Background(), domains.UserIDKey, userID.String())

		mockRepo.EXPECT().
			GetByID(gomock.Any(), userID).
			Return(nil, errors.New("db error"))

		_, err := uc.GetMe(ctx)
		assert.Error(t, err)
	})
}