package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Surname     null.String `json:"surname"`
	PhoneNumber null.String `json:"phone_number"`
}

type UserLoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequestDTO struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Surname  null.String `json:"surname"`
}

type UserResponseDTO struct {
	Token string `json:"token"`
}

type UserRepo struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      null.String
	PhoneNumber  null.String
	PasswordHash []byte
	Version      int
}
