package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
)

type ITokenator interface {
	CreateJWT(userID string, version int) (string, error)
	ParseJWT(tokenString string) (*jwt.JWTClaims, error)
}

//go:generate mockgen -source=user.go -destination=../repository/mocks/user_repository_mock.go -package=mocks IUserRepository
type IUserRepository interface {
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
	repo  IUserRepository
}

func NewAuthUsecase(repo IUserRepository, token ITokenator, log *logrus.Logger) *AuthUsecase {
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

	existed, err := u.repo.CheckUserExists(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if existed {
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

	token, err := u.token.CreateJWT(userRepo.ID.String(), userRepo.Version)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecase) Login(ctx context.Context, user models.UserLoginRequestDTO) (string, error) {
	userRepo, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
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

func (u *AuthUsecase) GetMe(ctx context.Context) (*models.User, error) {
	userIDStr, isExist := ctx.Value(utils.UserIDKey).(string)
	if !isExist {
		return nil, models.ErrUserNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, models.ErrInvalidUserID
	}

	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := userRepo.ConvertToUser()
	if user == nil {
		return nil, models.ErrUserNotFound
	}

	return user, nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
