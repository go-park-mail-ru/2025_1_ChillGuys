package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
)

type ITokenator interface {
	CreateJWT(userID string, version int) (string, error)
	ParseJWT(tokenString string) (*jwt.JWTClaims, error)
}

//go:generate mockgen -source=auth.go -destination=../../infrastructure/repository/postgres/mocks/auth_repository_mock.go -package=mocks IAuthRepository
type IAuthRepository interface {
	CreateUser(context.Context, models.UserDB) error
	GetUserByEmail(context.Context, string) (*models.UserDB, error)
	GetUserByID(context.Context, uuid.UUID) (*models.UserDB, error)
	IncrementUserVersion(context.Context, string) error
	GetUserCurrentVersion(context.Context, string) (int, error)
	CheckUserVersion(context.Context, string, int) bool
	CheckUserExists(context.Context, string) (bool, error)
}

type AuthUsecase struct {
	token ITokenator
	repo  IAuthRepository
}

func NewAuthUsecase(repo IAuthRepository, token ITokenator) *AuthUsecase {
	return &AuthUsecase{
		repo:  repo,
		token: token,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, user dto.UserRegisterRequestDTO) (string, uuid.UUID, error) {
	const op = "AuthUsecase.Register"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("email", user.Email)
	
	passwordHash, err := GeneratePasswordHash(user.Password)
	if err != nil {
		logger.WithError(err).Error("generate password hash")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	existed, err := u.repo.CheckUserExists(ctx, user.Email)
	if err != nil {
		logger.WithError(err).Error("check user existence")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	if existed {
		logger.Warn("user already exists")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
	}

	userID := uuid.New()
	userDB := models.UserDB{
		ID:           userID,
		Email:        user.Email,
		Name:         user.Name,
		Surname:      user.Surname,
		PasswordHash: passwordHash,
		UserVersion: models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    userID,
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	if err = u.repo.CreateUser(ctx, userDB); err != nil {
		logger.WithError(err).Error("create user in repository")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.UserVersion.Version)
	if err != nil {
		logger.WithError(err).Error("create JWT token")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return token, userDB.ID, nil
}

func (u *AuthUsecase) Login(ctx context.Context, user dto.UserLoginRequestDTO) (string, uuid.UUID, error) {
	const op = "AuthUsecase.Login"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("email", user.Email)

	userDB, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("user not found")
		} else {
			logger.WithError(err).Error("get user by email")
		}
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		logger.Warn("invalid credentials")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials)
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.UserVersion.Version)
	if err != nil {
		logger.WithError(err).Error("create JWT token")
		return "", uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return token, userDB.ID, nil
}

func (u *AuthUsecase) Logout(ctx context.Context) error {
	const op = "AuthUsecase.Logout"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userID, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		logger.Warn("user ID not found in context")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	logger = logger.WithField("user_id", userID)
	if err := u.repo.IncrementUserVersion(ctx, userID); err != nil {
		logger.WithError(err).Error("increment user version")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
