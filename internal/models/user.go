package models

import "github.com/google/uuid"

type UserDTO struct {
	ID      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Version string    `json:"version"`
}

type UserRepo struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
	Version      int
}
