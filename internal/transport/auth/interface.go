package auth

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
)

//go:generate mockgen -source=auth.go -destination=../../usecase/mocks/auth_usecase_mock.go -package=mocks IAuthUsecase
type IAuthUsecase interface {
	Register(context.Context, dto.UserRegisterRequestDTO) (string, error)
	Login(context.Context, dto.UserLoginRequestDTO) (string, error)
	Logout(context.Context) error
}
