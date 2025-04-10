package user

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
)

type ITokenator interface {
	CreateJWT(userID string, version int) (string, error)
	ParseJWT(tokenString string) (*jwt.JWTClaims, error)
}

//go:generate mockgen -source=user.go -destination=../../infrastructure/repository/postgres/mocks/user_repository_mock.go -package=mocks IProductRepository
type IUserRepository interface {
	Create(context.Context, dto.UserDB) error
	GetByEmail(context.Context, string) (*dto.UserDB, error)
	GetByID(context.Context, uuid.UUID) (*dto.UserDB, error)
	IncrementVersion(context.Context, string) error
	GetCurrentVersion(context.Context, string) (int, error)
	CheckVersion(context.Context, string, int) bool
	CheckExistence(context.Context, string) (bool, error)
	UpdateImageURL(context.Context, uuid.UUID, string) error
}

type AuthUsecase struct {
	token        ITokenator
	repo         IUserRepository
	minioService minio.Provider
}

func NewAuthUsecase(repo IUserRepository, token ITokenator, minioService minio.Provider) *AuthUsecase {
	return &AuthUsecase{
		repo:         repo,
		token:        token,
		minioService: minioService,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, user dto.UserRegisterRequestDTO) (string, error) {
	const op = "authUsecaseRegister"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	passwordHash, err := helpers.GeneratePasswordHash(user.Password)
	if err != nil {
		logger.WithError(err).Error("failed to generate password hash")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	existed, err := u.repo.CheckExistence(ctx, user.Email)
	if err != nil {
		logger.WithError(err).Error("failed to check user existence")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if existed {
		logger.Warn("user already exists")
		return "", fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
	}

	userDB := dto.NewFromRequest(user, passwordHash)

	if err = u.repo.Create(ctx, userDB); err != nil {
		logger.WithError(err).Error("failed to create user in repository")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.UserVersion.Version)
	if err != nil {
		logger.WithError(err).Error("failed to create JWT token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (u *AuthUsecase) Login(ctx context.Context, user dto.UserLoginRequestDTO) (string, error) {
	const op = "authUsecaseLogin"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	userDB, err := u.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		logger.WithError(err).Error("failed to get user by email")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		logger.WithError(err).Warn("invalid credentials provided")
		return "", fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.UserVersion.Version)
	if err != nil {
		logger.WithError(err).Error("failed to create JWT token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (u *AuthUsecase) Logout(ctx context.Context) error {
	const op = "authUsecaseLogout"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, isExist := ctx.Value(domains.UserIDKey).(string)
	if !isExist {
		logger.Error("user ID not found in context")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	if err := u.repo.IncrementVersion(ctx, userID); err != nil {
		logger.WithError(err).Error("increment user version")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *AuthUsecase) GetMe(ctx context.Context) (*models.User, error) {
	const op = "authUsecaseGetMe"
    logger := logctx.GetLogger(ctx).WithField("op", op)
	logger.Info("start")

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userDB, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("get user by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user := userDB.ConvertToUser()
	if user == nil {
		logger.WithError(err).Error("user not found after conversion")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	return user, nil
}

func (u *AuthUsecase) UploadAvatar(ctx context.Context, fileData minio.FileData) (string, error) {
	const op = "authUsecaseUploadAvatar"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	avatar, err := u.minioService.CreateOne(ctx, fileData)
	if err != nil {
		logger.WithError(err).Error("create avatar in storage")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = u.repo.UpdateImageURL(ctx, userID, avatar.URL); err != nil {
		logger.WithError(err).Error("update user avatar URL")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return avatar.URL, nil
}