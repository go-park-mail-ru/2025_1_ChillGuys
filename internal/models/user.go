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

func (u *UserDB) ConvertToUser() *User {
	if u == nil {
        return nil
    }
	return &User{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Surname:     u.Surname,
		PhoneNumber: u.PhoneNumber,
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

func (u *UserDB) IsVersionValid(version int) bool {
	return u.Version == version
}
