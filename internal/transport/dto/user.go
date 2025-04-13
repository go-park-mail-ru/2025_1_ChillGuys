package dto

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type UserDTO struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Surname     null.String `json:"surname" swaggertype:"primitive,string"`
	ImageURL    null.String `json:"imageURL" swaggertype:"primitive,string"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

type UpdateUserProfileRequestDTO struct {
	Name        null.String `json:"name,omitempty"`
	Surname     null.String `json:"surname,omitempty"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
}

type UpdateUserEmailDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserPasswordDTO struct {
	OldPassword string `json:"OldPassword"`
	NewPassword string `json:"NewPassword"`
}
