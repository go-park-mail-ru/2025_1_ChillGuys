package dto

import (
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func (u *UserRegisterRequestDTO) ConvertToGrpcRegisterReq() *gen.RegisterReq {
	var surname *wrapperspb.StringValue
	if u.Surname.Valid {
		surname = wrapperspb.String(u.Surname.String)
	}

	return &gen.RegisterReq{
		Email:    u.Email,
		Password: u.Password,
		Name:     u.Name,
		Surname:  surname,
	}
}

type UserResponseDTO struct {
	Token string `json:"token"`
}

type ErrorResponseDTO struct {
	Message string `json:"message"`
}
