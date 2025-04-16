package user

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//go:generate mockgen -source=user.go -destination=../../infrastructure/repository/postgres/mocks/user_repository_mock.go -package=mocks IUserRepository
type IUserRepository interface {
	GetUserByEmail(context.Context, string) (*models.UserDB, error)
	GetUserByID(context.Context, uuid.UUID) (*models.UserDB, error)
	UpdateUserImageURL(context.Context, uuid.UUID, string) error
	UpdateUserProfile(context.Context, uuid.UUID, models.UpdateUserDB) error
	UpdateUserEmail(context.Context, uuid.UUID, string) error
	UpdateUserPassword(context.Context, uuid.UUID, []byte) error
}

type UserUsecase struct {
	log          *logrus.Logger
	token        auth.ITokenator
	repo         IUserRepository
	minioService minio.Provider
}

func NewUserUsecase(repo IUserRepository, token auth.ITokenator, log *logrus.Logger, minioService minio.Provider) *UserUsecase {
	return &UserUsecase{
		repo:         repo,
		token:        token,
		minioService: minioService,
	}
}

func (u *UserUsecase) GetMe(ctx context.Context) (*dto.UserDTO, error) {
	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		return nil, errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := userRepo.ConvertToUser()

	if user == nil {
		return nil, errs.ErrNotFound
	}

	return &dto.UserDTO{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Surname:     user.Surname,
		ImageURL:    user.ImageURL,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (u *UserUsecase) UploadAvatar(ctx context.Context, fileData minio.FileData) (string, error) {
	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		return "", errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", errs.ErrInvalidID
	}

	avatar, err := u.minioService.CreateOne(ctx, fileData)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	if err = u.repo.UpdateUserImageURL(ctx, userID, avatar.URL); err != nil {
		return "", err
	}

	return avatar.URL, nil
}

func (u *UserUsecase) UpdateUserProfile(ctx context.Context, user dto.UpdateUserProfileRequestDTO) error {
	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		return errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errs.ErrInvalidID
	}

	currentUser, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	userDB := models.UpdateUserDB{}

	if user.Name.Valid && strings.TrimSpace(user.Name.String) != "" {
		userDB.Name = user.Name.String
	} else {
		userDB.Name = currentUser.Name
	}

	if user.Surname.Valid && strings.TrimSpace(user.Surname.String) != "" {
		userDB.Surname = user.Surname
	} else {
		userDB.Surname = currentUser.Surname
	}

	if user.PhoneNumber.Valid && strings.TrimSpace(user.PhoneNumber.String) != "" {
		userDB.PhoneNumber = user.PhoneNumber
	} else {
		userDB.PhoneNumber = currentUser.PhoneNumber
	}

	return u.repo.UpdateUserProfile(ctx, userID, userDB)
}

func (u *UserUsecase) UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return errs.ErrInvalidID
	}

	userDB, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		return errs.ErrInvalidCredentials
	}

	return u.repo.UpdateUserEmail(ctx, userID, user.Email)
}

func (u *UserUsecase) UpdateUserPassword(ctx context.Context, user dto.UpdateUserPasswordDTO) error {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return errs.ErrInvalidID
	}

	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(userRepo.PasswordHash, []byte(user.OldPassword)); err != nil {
		return errs.ErrInvalidCredentials
	}

	passwordHash, err := auth.GeneratePasswordHash(user.NewPassword)
	if err != nil {
		return err
	}

	return u.repo.UpdateUserPassword(ctx, userID, passwordHash)
}
