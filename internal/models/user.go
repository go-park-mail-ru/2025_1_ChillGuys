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
	ImageURL    null.String `json:"imageURL" swaggertype:"primitive,string"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

type UpdateUserDB struct {
	Name        string
	Surname     null.String
	PhoneNumber null.String
}

type UserDB struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      null.String
	ImageURL     null.String
	PhoneNumber  null.String
	PasswordHash []byte
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
		ImageURL:    u.ImageURL,
		PhoneNumber: u.PhoneNumber,
	}
}
