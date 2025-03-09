package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Surname     null.String `json:"surname" swaggertype:"primitive,string"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

func (repo *UserRepo) ConvertToUser() *User {
	return &User{
		ID:          repo.ID,
		Email:       repo.Email,
		Name:        repo.Name,
		Surname:     repo.Surname,
		PhoneNumber: repo.PhoneNumber,
	}
}

type UserLoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequestDTO struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Name     string      `json:"name"`
	Surname  null.String `json:"surname,omitempty" swaggertype:"primitive,string"`
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

func (repo *UserRepo) IsVersionValid(version int) bool {
	return repo.Version == version
}
