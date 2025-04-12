package models

import (
	"github.com/google/uuid"
	"time"
)

type UserVersionDB struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Version   int
	UpdatedAt time.Time
}
