package models

import (
	"time"
	
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Text      string    `json:"text"`
	Title     string    `json:"title"`
	IsRead    bool      `json:"is_read"`
	UpdatedAt time.Time `json:"updated_at"`
}