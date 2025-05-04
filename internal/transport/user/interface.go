package user

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
)

//go:generate mockgen -source=user.go -destination=../../usecase/mocks/user_usecase_mock.go -package=mocks IUserUsecase
type IUserUsecase interface {
	GetMe(context.Context) (*dto.UserDTO, string, error)
	UploadAvatar(context.Context, minio.FileData) (string, error)
	UpdateUserProfile(context.Context, dto.UpdateUserProfileRequestDTO) error
	UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error
	UpdateUserPassword(context.Context, dto.UpdateUserPasswordDTO) error
	BecomeSeller(ctx context.Context, req dto.UpdateRoleRequest) error
}