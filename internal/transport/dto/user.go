package dto

import "github.com/guregu/null"

type UpdateUserProfileRequestDTO struct {
	Name        null.String `json:"name,omitempty"`
	Surname     null.String `json:"surname,omitempty"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

type UpdateUserEmailRequestDTO struct {
	Name        null.String `json:"name,omitempty"`
	Surname     null.String `json:"surname,omitempty"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

type UpdateUserDB struct {
	Name        string
	Surname     null.String
	PhoneNumber null.String
}

type UpdateUserEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserPassword struct {
	OldPassword string `json:"OldPassword"`
	NewPassword string `json:"NewPassword"`
}
