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

func (ur *UserDB) ConvertToUser() *User {
	return &User{
		ID:          ur.ID,
		Email:       ur.Email,
		Name:        ur.Name,
		Surname:     ur.Surname,
		PhoneNumber: ur.PhoneNumber,
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

type UserDB struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      null.String
	PhoneNumber  null.String
	PasswordHash []byte
	Version      int
}

func (ur *UserDB) IsVersionValid(version int) bool {
	return ur.Version == version
}
