package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/notification"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAllByUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())
	ctx = logctx.WithLogger(ctx, logrus.NewEntry(logrus.New()))

	now := time.Now()
	notificationsDB := []models.Notification{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "Test 1",
			Title:     "Test",
			IsRead:    false,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "Test 2",
			Title:     "Test",
			IsRead:    true,
			UpdatedAt: now.Add(-time.Hour),
		},
	}

	expectedResponse := dto.NotificationsListResponse{
		Notifications: []dto.NotificationResponse{
			{
				ID:        notificationsDB[0].ID,
				Text:      notificationsDB[0].Text,
				Title:     notificationsDB[0].Title,
				IsRead:    notificationsDB[0].IsRead,
				UpdatedAt: notificationsDB[0].UpdatedAt,
			},
			{
				ID:        notificationsDB[1].ID,
				Text:      notificationsDB[1].Text,
				Title:     notificationsDB[1].Title,
				IsRead:    notificationsDB[1].IsRead,
				UpdatedAt: notificationsDB[1].UpdatedAt,
			},
		},
		Total:       2,
		UnreadCount: 1,
	}

	mockRepo.EXPECT().GetAllByUser(ctx, userID, 0).Return(notificationsDB, nil)
	mockRepo.EXPECT().GetUnreadCount(ctx, userID).Return(1, nil)

	result, err := uc.GetAllByUser(ctx, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
}

func TestGetAllByUser_GetUserIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	_, err := uc.GetAllByUser(ctx, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestGetAllByUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())
	ctx = logctx.WithLogger(ctx, logrus.NewEntry(logrus.New()))

	mockRepo.EXPECT().GetAllByUser(ctx, userID, 0).Return(nil, errors.New("database error"))

	_, err := uc.GetAllByUser(ctx, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestGetAllByUser_UnreadCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())
	ctx = logctx.WithLogger(ctx, logrus.NewEntry(logrus.New()))

	notificationsDB := []models.Notification{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Text:      "Test",
			Title:     "Test",
			IsRead:    false,
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.EXPECT().GetAllByUser(ctx, userID, 0).Return(notificationsDB, nil)
	mockRepo.EXPECT().GetUnreadCount(ctx, userID).Return(0, errors.New("count error"))

	_, err := uc.GetAllByUser(ctx, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "count error")
}

func TestGetUnreadCount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())
	ctx = logctx.WithLogger(ctx, logrus.NewEntry(logrus.New()))

	expectedCount := 3

	mockRepo.EXPECT().GetUnreadCount(ctx, userID).Return(expectedCount, nil)

	count, err := uc.GetUnreadCount(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestGetUnreadCount_GetUserIDError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	_, err := uc.GetUnreadCount(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestGetUnreadCount_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), domains.UserIDKey{}, userID.String())
	ctx = logctx.WithLogger(ctx, logrus.NewEntry(logrus.New()))

	mockRepo.EXPECT().GetUnreadCount(ctx, userID).Return(0, errors.New("database error"))

	_, err := uc.GetUnreadCount(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestMarkAsRead_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	notificationID := uuid.New()

	mockRepo.EXPECT().UpdateReadStatus(ctx, notificationID, true).Return(nil)

	err := uc.MarkAsRead(ctx, notificationID)

	assert.NoError(t, err)
}

func TestMarkAsRead_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockINotificationRepository(ctrl)
	uc := notification.NewNotificationUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	notificationID := uuid.New()

	mockRepo.EXPECT().UpdateReadStatus(ctx, notificationID, true).Return(errors.New("database error"))

	err := uc.MarkAsRead(ctx, notificationID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}
