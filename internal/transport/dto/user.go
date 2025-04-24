package dto

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/wrapperspb"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
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
	Name        null.String `json:"name,omitempty" swaggertype:"primitive,string"`
	Surname     null.String `json:"surname,omitempty" swaggertype:"primitive,string"`
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

func (u *UserDTO) ConvertToGrpcUser() *gen.User {
	var surname *wrapperspb.StringValue
	if u.Surname.Valid {
		surname = wrapperspb.String(u.Surname.String)
	}

	var imageURL *wrapperspb.StringValue
	if u.ImageURL.Valid {
		imageURL = wrapperspb.String(u.ImageURL.String)
	}

	var phoneNumber *wrapperspb.StringValue
	if u.PhoneNumber.Valid {
		phoneNumber = wrapperspb.String(u.PhoneNumber.String)
	}

	return &gen.User{
		Id:          u.ID.String(),
		Email:       u.Email,
		Name:        u.Name,
		Surname:     surname,
		ImageURL:    imageURL,
		PhoneNumber: phoneNumber,
	}
}

func ConvertGrpcToUserDTO(u *gen.User) (*UserDTO, error) {
	id, err := uuid.Parse(u.Id)
	if err != nil {
		return nil, err
	}

	var surname null.String
	if u.Surname != nil {
		surname = null.StringFrom(u.Surname.Value)
	}

	var imageURL null.String
	if u.ImageURL != nil {
		imageURL = null.StringFrom(u.ImageURL.Value)
	}

	var phoneNumber null.String
	if u.PhoneNumber != nil {
		phoneNumber = null.StringFrom(u.PhoneNumber.Value)
	}

	return &UserDTO{
		ID:          id,
		Email:       u.Email,
		Name:        u.Name,
		Surname:     surname,
		ImageURL:    imageURL,
		PhoneNumber: phoneNumber,
	}, nil
}

func (u *UpdateUserProfileRequestDTO) ConvertToGrpcUpdateProfileReq() *gen.UpdateUserProfileRequest {
	var name *wrapperspb.StringValue
	if u.Name.Valid {
		name = wrapperspb.String(u.Name.String)
	}

	var surname *wrapperspb.StringValue
	if u.Surname.Valid {
		surname = wrapperspb.String(u.Surname.String)
	}

	var phoneNumber *wrapperspb.StringValue
	if u.PhoneNumber.Valid {
		phoneNumber = wrapperspb.String(u.PhoneNumber.String)
	}

	return &gen.UpdateUserProfileRequest{
		Name:        name,
		Surname:     surname,
		PhoneNumber: phoneNumber,
	}
}