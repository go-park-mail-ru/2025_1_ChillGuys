package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	log   *logrus.Logger
	token transport.ITokenator
	repo  transport.IUserRepository
}

func NewAuthUsecase(repo transport.IUserRepository, token transport.ITokenator, log *logrus.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo:  repo,
		token: token,
		log:   log,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, user models.UserRegisterRequestDTO) (string, error) {
	passwordHash, err := GeneratePasswordHash(user.Password)
	if err != nil {
		return "", err
	}

	existedUser, _ := u.repo.GetUserByEmail(ctx, user.Email)
	if existedUser != nil {
		return "", models.ErrUserAlreadyExists
	}

	userRepo := models.UserDB{
		ID:           uuid.New(),
		Email:        user.Email,
		Name:         user.Name,
		Surname:      user.Surname,
		PasswordHash: passwordHash,
		Version:      1,
	}

	if err = u.repo.CreateUser(ctx, userRepo); err != nil {
		return "", err
	}

	token, err := u.token.CreateJWT(userRepo.ID.String(), 1)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecase) Login(ctx context.Context, user models.UserLoginRequestDTO) (string, error) {
	userRepo, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return "", models.ErrUserNotFound
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(userRepo.PasswordHash, []byte(user.Password)); err != nil {
		return "", models.ErrInvalidCredentials
	}

	token, err := u.token.CreateJWT(userRepo.ID.String(), userRepo.Version)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecase) Logout(ctx context.Context) error {
	userID, isExist := ctx.Value(utils.UserIDKey).(string)
	if !isExist {
		return models.ErrUserNotFound
	}

	if err := u.repo.IncrementUserVersion(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) GetMe(ctx context.Context) (models.User, error) {
	userIDStr, isExist := ctx.Value(utils.UserIDKey).(string)
	if !isExist {
		return models.User{}, models.ErrUserNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return models.User{}, models.ErrInvalidUserID
	}

	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return *userRepo.ConvertToUser(), nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
