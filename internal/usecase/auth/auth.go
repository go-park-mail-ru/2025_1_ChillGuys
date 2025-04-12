package auth

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	log   *logrus.Logger
	token ITokenator
	repo  IAuthRepository
}

func NewAuthUsecase(repo IAuthRepository, token ITokenator, log *logrus.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo:  repo,
		token: token,
		log:   log,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, user dto.UserRegisterRequestDTO) (string, error) {
	passwordHash, err := GeneratePasswordHash(user.Password)
	if err != nil {
		return "", err
	}

	existed, err := u.repo.CheckUserExists(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if existed {
		return "", errs.ErrAlreadyExists
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
		return "", err
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.UserVersion.Version)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecase) Login(ctx context.Context, user dto.UserLoginRequestDTO) (string, error) {
	userDB, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		return "", errs.ErrInvalidCredentials
	}

	token, err := u.token.CreateJWT(userDB.ID.String(), userDB.UserVersion.Version)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecase) Logout(ctx context.Context) error {
	userID, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist {
		return errs.ErrNotFound
	}

	if err := u.repo.IncrementUserVersion(ctx, userID); err != nil {
		return err
	}

	return nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
