package dto

import (
	"time"
	
	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	Title     string    `json:"title"`
	IsRead    bool      `json:"is_read"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotificationsListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int                   `json:"total"`
	UnreadCount   int                   `json:"unread_count"`
}

type CreateNotificationRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Text   string    `json:"text" validate:"required"`
	Title  string    `json:"title" validate:"required"`
}

type UpdateNotificationStatusRequest struct {
	ID     uuid.UUID `json:"id"`
	IsRead bool      `json:"is_read"`
}