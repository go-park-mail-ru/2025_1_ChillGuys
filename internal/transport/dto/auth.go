package dto

import (
	"github.com/guregu/null"
)

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

type ErrorResponseDTO struct {
	Message string `json:"message"`
}
