package tests

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository/mocks"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	user := models.UserDB{
		ID:      uuid.New(),
		Email:   "test@example.com",
		Version: 1,
	}

	mockRepo.EXPECT().CreateUser(context.Background(), user).Return(nil)
	mockRepo.EXPECT().GetUserByEmail(context.Background(), user.Email).Return(&user, nil)

	err := mockRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	storedUser, err := mockRepo.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, storedUser.ID)
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	user := models.UserDB{
		ID:      uuid.New(),
		Email:   "test@example.com",
		Version: 1,
	}

	mockRepo.EXPECT().CreateUser(context.Background(), user).Return(nil)
	mockRepo.EXPECT().GetUserByEmail(context.Background(), user.Email).Return(&user, nil)
	mockRepo.EXPECT().GetUserByEmail(context.Background(), "nonexistent@example.com").Return(nil, models.ErrUserNotFound)

	err := mockRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	storedUser, err := mockRepo.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, storedUser.ID)

	_, err = mockRepo.GetUserByEmail(context.Background(), "nonexistent@example.com")
	assert.Equal(t, models.ErrUserNotFound, err)
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)

	userID := uuid.Must(uuid.Parse("33b92ff8-4bb9-4334-afb5-3a7408ed3ec5"))
	nonExistentID := uuid.Must(uuid.Parse("79b09acf-ca4d-496e-8dfe-7bcb7ca5d868"))

	user := models.UserDB{
		ID:      userID,
		Email:   "test@example.com",
		Version: 1,
	}

	mockRepo.EXPECT().CreateUser(context.Background(), user).Return(nil)
	mockRepo.EXPECT().GetUserByID(context.Background(), userID).Return(&user, nil)
	mockRepo.EXPECT().GetUserByID(context.Background(), nonExistentID).Return(nil, models.ErrUserNotFound)

	err := mockRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	storedUser, err := mockRepo.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, storedUser.ID)

	_, err = mockRepo.GetUserByID(context.Background(), nonExistentID)
	assert.Equal(t, models.ErrUserNotFound, err)
}

func TestIncrementUserVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	user := models.UserDB{
		ID:      uuid.New(),
		Email:   "test@example.com",
		Version: 1,
	}

	mockRepo.EXPECT().CreateUser(context.Background(), user).Return(nil)
	mockRepo.EXPECT().GetUserCurrentVersion(context.Background(), user.ID.String()).Return(1, nil)
	mockRepo.EXPECT().IncrementUserVersion(context.Background(), user.ID.String()).Return(nil)
	mockRepo.EXPECT().GetUserCurrentVersion(context.Background(), user.ID.String()).Return(2, nil)
	mockRepo.EXPECT().IncrementUserVersion(context.Background(), "nonexistent").Return(models.ErrUserNotFound)

	err := mockRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	version, err := mockRepo.GetUserCurrentVersion(context.Background(), user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 1, version)

	err = mockRepo.IncrementUserVersion(context.Background(), user.ID.String())
	assert.NoError(t, err)

	version, err = mockRepo.GetUserCurrentVersion(context.Background(), user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 2, version)

	err = mockRepo.IncrementUserVersion(context.Background(), "nonexistent")
	assert.Equal(t, models.ErrUserNotFound, err)
}

func TestCheckUserVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	user := models.UserDB{
		ID:      uuid.New(),
		Email:   "test@example.com",
		Version: 1,
	}

	mockRepo.EXPECT().CreateUser(context.Background(), user).Return(nil)
	mockRepo.EXPECT().CheckUserVersion(context.Background(), user.ID.String(), 1).Return(true)
	mockRepo.EXPECT().CheckUserVersion(context.Background(), user.ID.String(), 2).Return(false)
	mockRepo.EXPECT().CheckUserVersion(context.Background(), "nonexistent", 1).Return(false)

	err := mockRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	isValid := mockRepo.CheckUserVersion(context.Background(), user.ID.String(), 1)
	assert.True(t, isValid)

	isValid = mockRepo.CheckUserVersion(context.Background(), user.ID.String(), 2)
	assert.False(t, isValid)

	isValid = mockRepo.CheckUserVersion(context.Background(), "nonexistent", 1)
	assert.False(t, isValid)
}
