package tests

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetMe_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)

	userID := uuid.New()
	userDB := &models.UserDB{
		ID:          userID,
		Email:       "test@example.com",
		Name:        "Test",
		Surname:     null.NewString("User", true),
		ImageURL:    null.NewString("http://example.com/avatar.jpg", true),
		PhoneNumber: null.NewString("1234567890", true),
	}

	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())
	mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userDB, nil)

	usecase := user.NewUserUsecase(mockRepo, nil, logrus.New(), nil)

	userDTO, err := usecase.GetMe(ctx)
	require.NoError(t, err)
	require.Equal(t, userDB.ID, userDTO.ID)
	require.Equal(t, userDB.Email, userDTO.Email)
}

func TestUpdateUserProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())

	currentUser := &models.UserDB{
		Name:        "OldName",
		Surname:     null.NewString("OldSurname", true),
		PhoneNumber: null.NewString("111111", true),
	}
	mockRepo.EXPECT().GetUserByID(ctx, userID).Return(currentUser, nil)

	update := dto.UpdateUserProfileRequestDTO{
		Name:        null.NewString("NewName", true),
		Surname:     null.NewString("", false),
		PhoneNumber: null.NewString("222222", true),
	}

	expected := models.UpdateUserDB{
		Name:        "NewName",
		Surname:     currentUser.Surname,
		PhoneNumber: update.PhoneNumber,
	}

	mockRepo.EXPECT().UpdateUserProfile(ctx, userID, expected).Return(nil)

	usecase := user.NewUserUsecase(mockRepo, nil, logrus.New(), nil)

	err := usecase.UpdateUserProfile(ctx, update)
	require.NoError(t, err)
}

func TestUpdateUserEmail_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())

	password := "password123"
	hashed, _ := auth.GeneratePasswordHash(password)

	userDB := &models.UserDB{
		PasswordHash: hashed,
	}

	mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userDB, nil)
	mockRepo.EXPECT().UpdateUserEmail(ctx, userID, "new@example.com").Return(nil)

	usecase := user.NewUserUsecase(mockRepo, nil, logrus.New(), nil)

	dto := dto.UpdateUserEmailDTO{
		Email:    "new@example.com",
		Password: password,
	}

	err := usecase.UpdateUserEmail(ctx, dto)
	require.NoError(t, err)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())

	oldPassword := "oldpass"
	newPassword := "newpass"

	oldHash, _ := auth.GeneratePasswordHash(oldPassword)
	us := &models.UserDB{PasswordHash: oldHash}

	mockRepo.EXPECT().GetUserByID(ctx, userID).Return(us, nil)
	mockRepo.EXPECT().UpdateUserPassword(ctx, userID, gomock.Any()).Return(nil)

	usecase := user.NewUserUsecase(mockRepo, nil, logrus.New(), nil)

	dto := dto.UpdateUserPasswordDTO{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	err := usecase.UpdateUserPassword(ctx, dto)
	require.NoError(t, err)
}
