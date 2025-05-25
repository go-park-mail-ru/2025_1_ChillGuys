package notification

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/notification"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/google/uuid"
)

type INotificationUsecase interface {
	GetAllByUser(ctx context.Context, offset int) (dto.NotificationsListResponse, error)
	GetUnreadCount(ctx context.Context) (int, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}

type NotificationUsecase struct {
	repo notification.INotificationRepository
}

func NewNotificationUsecase(repo notification.INotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{repo: repo}
}

func (u *NotificationUsecase) GetAllByUser(ctx context.Context, offset int) (dto.NotificationsListResponse, error) {
	const op = "NotificationUsecase.GetAllByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return dto.NotificationsListResponse{}, fmt.Errorf("%s: %w", op, err)
    }
	
	logger = logger.WithField("user_id", userID)
	notificationsDB, err := u.repo.GetAllByUser(ctx, userID, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get notifications")
		return dto.NotificationsListResponse{}, err
	}

	notifications := make([]dto.NotificationResponse, 0, len(notificationsDB))
	for _, n := range notificationsDB {
		notifications = append(notifications, dto.NotificationResponse{
			ID:        n.ID,
			Text:      n.Text,
			Title:     n.Title,
			IsRead:    n.IsRead,
			UpdatedAt: n.UpdatedAt,
		})
	}

	count, err := u.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get count unrad")
        return dto.NotificationsListResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.NotificationsListResponse{
		Notifications: notifications,
		Total:         len(notifications),
		UnreadCount: count,
	}, nil
}

func (u *NotificationUsecase) GetUnreadCount(ctx context.Context) (int, error) {
	const op = "NotificationUsecase.GetAllByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return 0, fmt.Errorf("%s: %w", op, err)
    }
	
	logger = logger.WithField("user_id", userID)

	count, err := u.repo.GetUnreadCount(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("failed to get count unrad")
        return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (u *NotificationUsecase) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	const op = "NotificationUsecase.GetAllByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	err := u.repo.UpdateReadStatus(ctx, id, true)
	if err != nil {
		logger.WithError(err).Error("failed to mark read")
        return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}