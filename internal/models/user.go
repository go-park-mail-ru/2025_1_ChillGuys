package models

import "github.com/google/uuid"

type UserDTO struct {
	ID      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Version string    `json:"version"`
}

type UserRepo struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      string
	PasswordHash []byte
	Version      int
}
