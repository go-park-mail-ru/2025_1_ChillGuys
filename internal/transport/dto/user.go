package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type UserDB struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      null.String
	ImageURL     null.String
	PhoneNumber  null.String
	PasswordHash []byte
	UserVersion  models.UserVersionDB
}

func (u *UserDB) ConvertToUser() *models.User {
	if u == nil {
		return nil
	}
	return &models.User{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Surname:     u.Surname,
		ImageURL:    u.ImageURL,
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

type UpdateUserProfileRequestDTO struct {
	Email       null.String `json:"email,omitempty"`
	Name        null.String `json:"name,omitempty"`
	Surname     null.String `json:"surname,omitempty"`
	Password    null.String `json:"password,omitempty"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

type UpdateUserDB struct {
	Email        string
	Name         string
	Surname      null.String
	ImageURL     null.String
	PhoneNumber  null.String
	PasswordHash []byte
}

type UserResponseDTO struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (u *UserDB) IsVersionValid(version int) bool {
	return u.UserVersion.Version == version
}
