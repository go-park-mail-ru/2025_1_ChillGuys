package models

import "github.com/google/uuid"

type Category struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
}