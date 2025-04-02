package user

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
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

// FIXME: переписать путь моков
//
//go:generate mockgen -source=user.go -destination=./mocks/user_repository_mock.go -package=mocks IUserRepository
type IUserRepository interface {
	CreateUser(context.Context, models.UserDB) error
	GetUserByEmail(context.Context, string) (*models.UserDB, error)
	GetUserByID(context.Context, uuid.UUID) (*models.UserDB, error)
	IncrementUserVersion(context.Context, string) error
	GetUserCurrentVersion(context.Context, string) (int, error)
	CheckUserVersion(context.Context, string, int) bool
	CheckUserExists(context.Context, string) (bool, error)
	UpdateUserImageURL(context.Context, uuid.UUID, string) error
}

type AuthUsecase struct {
	log          *logrus.Logger
	token        ITokenator
	repo         IUserRepository
	minioService minio.Client
}

func NewAuthUsecase(repo IUserRepository, token ITokenator, log *logrus.Logger, minioService minio.Client) *AuthUsecase {
	return &AuthUsecase{
		repo:         repo,
		token:        token,
		log:          log,
		minioService: minioService,
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
		return "", errs.ErrUserAlreadyExists
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

func (u *AuthUsecase) Login(ctx context.Context, user models.UserLoginRequestDTO) (string, error) {
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
	userID, isExist := ctx.Value(cookie.UserIDKey).(string)
	if !isExist {
		return errs.ErrUserNotFound
	}

	if err := u.repo.IncrementUserVersion(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) GetMe(ctx context.Context) (*models.User, error) {
	userIDStr, isExist := ctx.Value(cookie.UserIDKey).(string)
	if !isExist {
		return nil, errs.ErrUserNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errs.ErrInvalidUserID
	}

	userRepo, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := userRepo.ConvertToUser()
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	return user, nil
}

func (u *AuthUsecase) UploadAvatar(ctx context.Context, fileData minio.FileDataType) (string, error) {
	userIDStr, isExist := ctx.Value(cookie.UserIDKey).(string)
	if !isExist {
		return "", errs.ErrUserNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", errs.ErrInvalidUserID
	}

	avatar, err := u.minioService.CreateOne(ctx, fileData)
	if err != nil {
		return "", err
	}

	if err = u.repo.UpdateUserImageURL(ctx, userID, avatar.URL); err != nil {
		return "", err
	}

	return avatar.URL, nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
