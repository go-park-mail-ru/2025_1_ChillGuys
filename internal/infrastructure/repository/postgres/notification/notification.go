package notification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryCreateNotification = `
		INSERT INTO bazaar.notification 
		(id, user_id, text, title, is_read) 
		VALUES ($1, $2, $3, $4, $5)`

	queryGetAllNotifications = `
		SELECT id, user_id, text, title, is_read, updated_at 
		FROM bazaar.notification 
		WHERE user_id = $1 
		ORDER BY updated_at DESC`

	queryGetUnreadCount = `
		SELECT COUNT(*) 
		FROM bazaar.notification 
		WHERE user_id = $1 AND is_read = false`

	queryUpdateReadStatus = `
		UPDATE bazaar.notification 
		SET is_read = $1 
		WHERE id = $2`
)

type INotificationRepository interface {
	Create(ctx context.Context, notification models.Notification) error
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Notification, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	UpdateReadStatus(ctx context.Context, id uuid.UUID, isRead bool) error
}

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification models.Notification) error {
	const op = "NotificationRepository.Create"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	_, err := r.db.ExecContext(ctx, queryCreateNotification,
		notification.ID,
		notification.UserID,
		notification.Text,
		notification.Title,
		notification.IsRead,
	)

	if err != nil {
		logger.WithError(err).Error("create notification")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *NotificationRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	const op = "NotificationRepository.GetAllByUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	rows, err := r.db.QueryContext(ctx, queryGetAllNotifications, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("no notifications found")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("query notifications")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Text,
			&n.Title,
			&n.IsRead,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	const op = "NotificationRepository.GetUnreadCount"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	var count int
	err := r.db.QueryRowContext(ctx, queryGetUnreadCount, userID).Scan(&count)
	if err != nil {
		logger.WithError(err).Error(err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (r *NotificationRepository) UpdateReadStatus(ctx context.Context, id uuid.UUID, isRead bool) error {
	const op = "NotificationRepository.UpdateReadStatus"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	_, err := r.db.ExecContext(ctx, queryUpdateReadStatus, isRead, id)
	if err != nil {
		logger.WithError(err).Error(err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}