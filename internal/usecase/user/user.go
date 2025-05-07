package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=user.go -destination=../../infrastructure/repository/postgres/mocks/user_repository_mock.go -package=mocks IUserRepository
type IUserRepository interface {
	GetUserByEmail(context.Context, string) (*models.UserDB, error)
	GetUserByID(context.Context, uuid.UUID) (*models.UserDB, error)
	UpdateUserImageURL(context.Context, uuid.UUID, string) error
	UpdateUserProfile(context.Context, uuid.UUID, models.UpdateUserDB) error
	UpdateUserEmail(context.Context, uuid.UUID, string) error
	UpdateUserPassword(context.Context, uuid.UUID, []byte) error
	CreateSellerAndUpdateRole(ctx context.Context, userID uuid.UUID, title, description string)  error
}

type UserUsecase struct {
	token        auth.ITokenator
	repo         IUserRepository
	minioService minio.Provider
}

func NewUserUsecase(repo IUserRepository, token auth.ITokenator, minioService minio.Provider) *UserUsecase {
	return &UserUsecase{
		repo:         repo,
		token:        token,
		minioService: minioService,
	}
}

func (u *UserUsecase) GetMe(ctx context.Context) (*dto.UserDTO, string, error) {
	const op = "UserUsecase.GetMe"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Warn("user ID not found in context")
		return nil, "", fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("invalid user ID format")
		return nil, "", fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	role, isExist := ctx.Value(domains.RoleKey{}).(string)
	fmt.Println(role, isExist)
	if !isExist {
		logger.Warn("role not found in context")
		return nil, "", fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger = logger.WithField("user_id", userID)
	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("user not found")
			return nil, "", fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("get user from repository")
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	user := userRepo.ConvertToUser()
	if user == nil {
		logger.Error("failed to convert user from db model")
		return nil, "", fmt.Errorf("%s: %w", op, errs.ErrBusinessLogic)
	}

	userDTO := &dto.UserDTO{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Surname:     user.Surname,
		ImageURL:    user.ImageURL,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.String(),
	}

	if role != user.Role.String() {
		token, err := u.token.CreateJWT(userIDStr, user.Role.String())
		if err != nil {
			logger.WithError(err).Error("create JWT token")
			return nil, "", fmt.Errorf("%s: %w", op, err)
		}
		return userDTO, token, nil
	}

	return userDTO, "", nil
}

func (u *UserUsecase) UploadAvatar(ctx context.Context, fileData minio.FileData) (string, error) {
	const op = "UserUsecase.UploadAvatar"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Warn("user ID not found in context")
		return "", fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("invalid user ID format")
		return "", fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	logger = logger.WithField("user_id", userID)
	avatar, err := u.minioService.CreateOne(ctx, fileData)
	if err != nil {
		logger.WithError(err).Error("upload avatar to storage")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = u.repo.UpdateUserImageURL(ctx, userID, avatar.URL); err != nil {
		logger.WithError(err).Error("update user avatar URL")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return avatar.URL, nil
}

func (u *UserUsecase) UpdateUserProfile(ctx context.Context, user dto.UpdateUserProfileRequestDTO) error {
	const op = "UserUsecase.UpdateUserProfile"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Warn("user ID not found in context")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("invalid user ID format")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	logger = logger.WithField("user_id", userID)
	currentUser, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("user not found")
			return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("get current user data")
		return fmt.Errorf("%s: %w", op, err)
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

	if err := u.repo.UpdateUserProfile(ctx, userID, userDB); err != nil {
		logger.WithError(err).Error("update user profile")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserUsecase) UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error {
	const op = "UserUsecase.UpdateUserEmail"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	logger = logger.WithField("user_id", userID)
	userDB, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("user not found")
			return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("get user data")
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		logger.Warn("invalid password provided")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	if err := u.repo.UpdateUserEmail(ctx, userID, user.Email); err != nil {
		logger.WithError(err).Error("update user email")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserUsecase) UpdateUserPassword(ctx context.Context, user dto.UpdateUserPasswordDTO) error {
	const op = "UserUsecase.UpdateUserPassword"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}

	logger = logger.WithField("user_id", userID)
	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("user not found")
			return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		logger.WithError(err).Error("get user data")
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(userRepo.PasswordHash, []byte(user.OldPassword)); err != nil {
		logger.Warn("invalid old password provided")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	passwordHash, err := auth.GeneratePasswordHash(user.NewPassword)
	if err != nil {
		logger.WithError(err).Error("generate password hash")
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := u.repo.UpdateUserPassword(ctx, userID, passwordHash); err != nil {
		logger.WithError(err).Error("update user password")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserUsecase) BecomeSeller(ctx context.Context, req dto.UpdateRoleRequest) error {
    const op = "UserUsecase.BecomeSeller"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
    if !isExist {
        logger.Warn("user ID not found in context")
        return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        logger.WithError(err).Error("invalid user ID format")
        return fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
    }

    err = u.repo.CreateSellerAndUpdateRole(ctx, userID, req.Title, req.Description)
    if err != nil {
        logger.WithError(err).Error("failed to become seller")
        return fmt.Errorf("%s: %w", op, err)
    }

    return nil
}