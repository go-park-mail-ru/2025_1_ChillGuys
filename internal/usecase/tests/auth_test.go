package tests

//
//import (
//	"context"
//	"errors"
//	user2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
//	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
//	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
//	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
//	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"github.com/google/uuid"
//	"github.com/guregu/null"
//	"github.com/sirupsen/logrus"
//	"github.com/stretchr/testify/assert"
//	"golang.org/x/crypto/bcrypt"
//
//	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
//)
//
//func TestAuthUsecase_Register(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockRepo := user2.NewMockIAuthRepository(ctrl)
//	mockToken := user2.NewMockITokenator(ctrl)
//	logger := logrus.New()
//
//	uc := auth.NewAuthUsecase(mockRepo, mockToken, logger)
//
//	testUser := dto.UserRegisterRequestDTO{
//		Email:    "test@example.com",
//		Password: "password123",
//		Name:     "Test",
//		Surname:  null.StringFrom("User"),
//	}
//
//	t.Run("Success", func(t *testing.T) {
//		mockRepo.EXPECT().
//			CheckUserExists(gomock.Any(), testUser.Email).
//			Return(false, nil).
//			Times(1)
//
//		mockRepo.EXPECT().
//			CreateUser(gomock.Any(), gomock.All(
//				gomock.AssignableToTypeOf(models.UserDB{}),
//				gomock.Not(gomock.Nil()),
//			)).
//			Return(nil).
//			Times(1)
//
//		mockToken.EXPECT().
//			CreateJWT(gomock.Any(), gomock.Any()).
//			Return("test-token", nil).
//			Times(1)
//
//		token, err := uc.Register(context.Background(), testUser)
//
//		assert.NoError(t, err)
//		assert.Equal(t, "test-token", token)
//	})
//
//	t.Run("CheckUserExistsError", func(t *testing.T) {
//		mockRepo.EXPECT().
//			CheckUserExists(gomock.Any(), testUser.Email).
//			Return(false, errors.New("db error")).
//			Times(1)
//
//		token, err := uc.Register(context.Background(), testUser)
//
//		assert.Error(t, err)
//		assert.Empty(t, token)
//	})
//
//	t.Run("CreateUserError", func(t *testing.T) {
//		mockRepo.EXPECT().
//			CheckUserExists(gomock.Any(), testUser.Email).
//			Return(false, nil).
//			Times(1)
//
//		mockRepo.EXPECT().
//			CreateUser(gomock.Any(), gomock.Any()).
//			Return(errors.New("db error")).
//			Times(1)
//
//		token, err := uc.Register(context.Background(), testUser)
//
//		assert.Error(t, err)
//		assert.Empty(t, token)
//	})
//
//	t.Run("CreateJWTError", func(t *testing.T) {
//		mockRepo.EXPECT().
//			CheckUserExists(gomock.Any(), testUser.Email).
//			Return(false, nil).
//			Times(1)
//
//		mockRepo.EXPECT().
//			CreateUser(gomock.Any(), gomock.Any()).
//			Return(nil).
//			Times(1)
//
//		mockToken.EXPECT().
//			CreateJWT(gomock.Any(), 1).
//			Return("", errors.New("jwt error")).
//			Times(1)
//
//		token, err := uc.Register(context.Background(), testUser)
//
//		assert.Error(t, err)
//		assert.Empty(t, token)
//	})
//}
//
//func TestAuthUsecase_Login(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockRepo := user2.NewMockIAuthRepository(ctrl)
//	mockToken := user2.NewMockITokenator(ctrl)
//	logger := logrus.New()
//
//	uc := auth.NewAuthUsecase(mockRepo, mockToken, logger)
//
//	testUserID := uuid.New()
//
//	testUser := dto.UserLoginRequestDTO{
//		Email:    "test@example.com",
//		Password: "password123",
//	}
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.MinCost)
//	testUserDB := &models.UserDB{
//		ID:           uuid.New(),
//		Email:        testUser.Email,
//		PasswordHash: hashedPassword,
//	}
//
//	t.Run("Success", func(t *testing.T) {
//		userVersion := models.UserVersionDB{
//			ID:        uuid.New(),
//			UserID:    testUserID,
//			Version:   1,
//			UpdatedAt: time.Now(),
//		}
//
//		userDB := &models.UserDB{
//			ID:           testUserID,
//			Email:        testUser.Email,
//			PasswordHash: hashedPassword,
//			UserVersion:  userVersion,
//		}
//
//		mockRepo.EXPECT().
//			GetUserByEmail(gomock.Any(), testUser.Email).
//			Return(userDB, nil).
//			Times(1)
//
//		mockToken.EXPECT().
//			CreateJWT(userDB.ID.String(), userDB.UserVersion.Version).
//			Return("test-token", nil).
//			Times(1)
//
//		token, err := uc.Login(context.Background(), testUser)
//
//		assert.NoError(t, err)
//		assert.Equal(t, "test-token", token)
//	})
//
//	t.Run("UserNotFound", func(t *testing.T) {
//		mockRepo.EXPECT().
//			GetUserByEmail(gomock.Any(), testUser.Email).
//			Return(nil, errors.New("not found")).
//			Times(1)
//
//		token, err := uc.Login(context.Background(), testUser)
//
//		assert.Error(t, err)
//		assert.Empty(t, token)
//	})
//
//	t.Run("InvalidCredentials", func(t *testing.T) {
//		invalidUserDB := *testUserDB
//		invalidUserDB.PasswordHash = []byte("wrong-hash")
//
//		mockRepo.EXPECT().
//			GetUserByEmail(gomock.Any(), testUser.Email).
//			Return(&invalidUserDB, nil).
//			Times(1)
//
//		token, err := uc.Login(context.Background(), testUser)
//
//		assert.Error(t, err)
//		assert.Equal(t, errs.ErrInvalidCredentials, err)
//		assert.Empty(t, token)
//	})
//
//	t.Run("CreateJWTError", func(t *testing.T) {
//		mockRepo.EXPECT().
//			GetUserByEmail(gomock.Any(), testUser.Email).
//			Return(testUserDB, nil).
//			Times(1)
//
//		mockToken.EXPECT().
//			CreateJWT(testUserDB.ID.String(), testUserDB.UserVersion.Version).
//			Return("", errors.New("jwt error")).
//			Times(1)
//
//		token, err := uc.Login(context.Background(), testUser)
//
//		assert.Error(t, err)
//		assert.Empty(t, token)
//	})
//
//}
//
//func TestAuthUsecase_Logout(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockRepo := user2.NewMockIAuthRepository(ctrl)
//	mockToken := user2.NewMockITokenator(ctrl)
//	logger := logrus.New()
//
//	uc := auth.NewAuthUsecase(mockRepo, mockToken, logger)
//
//	testUserID := uuid.New().String()
//
//	t.Run("Success", func(t *testing.T) {
//		ctx := context.WithValue(context.Background(), domains.UserIDKey{}, testUserID)
//
//		mockRepo.EXPECT().
//			IncrementUserVersion(gomock.Any(), testUserID).
//			Return(nil).
//			Times(1)
//
//		err := uc.Logout(ctx)
//
//		assert.NoError(t, err)
//	})
//
//	t.Run("UserNotFoundInContext", func(t *testing.T) {
//		err := uc.Logout(context.Background())
//
//		assert.Error(t, err)
//		assert.Equal(t, errs.ErrNotFound, err)
//	})
//
//	t.Run("IncrementVersionError", func(t *testing.T) {
//		ctx := context.WithValue(context.Background(), domains.UserIDKey{}, testUserID)
//
//		mockRepo.EXPECT().
//			IncrementUserVersion(gomock.Any(), testUserID).
//			Return(errors.New("db error")).
//			Times(1)
//
//		err := uc.Logout(ctx)
//
//		assert.Error(t, err)
//	})
//}
