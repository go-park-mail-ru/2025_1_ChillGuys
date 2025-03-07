package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	PhoneNumber string    `json:"phone_number"`
}

type UserLoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type UserResponseDTO struct {
	Token string `json:"token"`
}

type UserRepo struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      string
	PhoneNumber  string
	PasswordHash []byte
	Version      int
}
