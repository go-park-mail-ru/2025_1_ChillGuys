package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
)

type ITokenator interface {
	CreateJWT(userID string, role string) (string, error)
	ParseJWT(tokenString string) (*jwt.JWTClaims, error)
}

//go:generate mockgen -source=auth.go -destination=../../infrastructure/repository/postgres/mocks/auth_repository_mock.go -package=mocks IAuthRepository
type IAuthRepository interface {
	CreateUser(context.Context, models.UserDB) error
	GetUserByEmail(context.Context, string) (*models.UserDB, error)
	GetUserByID(context.Context, uuid.UUID) (*models.UserDB, error)
	CheckUserExists(context.Context, string) (bool, error)
}

type IAuthRedisRepository interface {
    AddToBlacklist(ctx context.Context, userID, token string) error
    IsInBlacklist(ctx context.Context, userID, token string) (bool, error)
}

type AuthUsecase struct {
	token     ITokenator
	repo      IAuthRepository
	redisRepo IAuthRedisRepository 
}

func NewAuthUsecase(repo IAuthRepository, redisRepo IAuthRedisRepository , token ITokenator) *AuthUsecase {
	return &AuthUsecase{
		repo:      repo,
		redisRepo: redisRepo,
		token:     token,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, user dto.UserRegisterRequestDTO) (string, error) {
	const op = "AuthUsecase.Register"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("email", user.Email)

	passwordHash, err := GeneratePasswordHash(user.Password)
	if err != nil {
		logger.WithError(err).Error("generate password hash")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	existed, err := u.repo.CheckUserExists(ctx, user.Email)
	if err != nil {
		logger.WithError(err).Error("check user existence")
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if existed {
		logger.Warn("user already exists")
		return "", fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
	}

	userID := uuid.New()
	userDB := models.UserDB{
		ID:           userID,
		Email:        user.Email,
		Name:         user.Name,
		Surname:      user.Surname,
		PasswordHash: passwordHash,
		Role:         models.RoleBuyer,
	}

	if err = u.repo.CreateUser(ctx, userDB); err != nil {
		logger.WithError(err).Error("create user in repository")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.Role.String())
	if err != nil {
		logger.WithError(err).Error("create JWT token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (u *AuthUsecase) Login(ctx context.Context, user dto.UserLoginRequestDTO) (string, error) {
	const op = "AuthUsecase.Login"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("email", user.Email)

	userDB, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("user not found")
		} else {
			logger.WithError(err).Error("get user by email")
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		logger.Warn("invalid credentials")
		return "", fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.Role.String())
	if err != nil {
		logger.WithError(err).Error("create JWT token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, token string) error {
	const op = "AuthUsecase.Logout"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	claims, err := u.token.ParseJWT(token)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidToken)
	}

	// Add token to blacklist with userID association
	if err := u.redisRepo.AddToBlacklist(ctx, claims.UserID, token); err != nil {
		logger.WithError(err).Error("failed to add token to blacklist")
		return fmt.Errorf("%s: %w", op, errs.ErrInternal)
	}

	return nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
