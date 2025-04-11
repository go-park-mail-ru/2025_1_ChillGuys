package user

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"strings"
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

//go:generate mockgen -source=user.go -destination=../../infrastructure/repository/postgres/mocks/user_repository_mock.go -package=mocks IProductRepository
type IUserRepository interface {
	CreateUser(context.Context, dto.UserDB) error
	GetUserByEmail(context.Context, string) (*dto.UserDB, error)
	GetUserByID(context.Context, uuid.UUID) (*dto.UserDB, error)
	IncrementUserVersion(context.Context, string) error
	GetUserCurrentVersion(context.Context, string) (int, error)
	CheckUserVersion(context.Context, string, int) bool
	CheckUserExists(context.Context, string) (bool, error)
	UpdateUserImageURL(context.Context, uuid.UUID, string) error
	UpdateUserProfile(context.Context, uuid.UUID, dto.UpdateUserDB) error
	UpdateUserEmail(context.Context, uuid.UUID, string) error
	UpdateUserPassword(context.Context, uuid.UUID, []byte) error
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
	userDB := dto.UserDB{
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
	userID, isExist := ctx.Value(domains.UserIDKey).(string)
	if !isExist {
		return errs.ErrNotFound
	}

	if err := u.repo.IncrementUserVersion(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) GetMe(ctx context.Context) (*models.User, error) {
	userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
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

	return user, nil
}

func (u *AuthUsecase) UploadAvatar(ctx context.Context, fileData minio.FileDataType) (string, error) {
	userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
	if !isExist {
		return "", errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", errs.ErrInvalidID
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

func (u *AuthUsecase) UpdateUserProfile(ctx context.Context, user dto.UpdateUserProfileRequestDTO) error {
	userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
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

	userDB := dto.UpdateUserDB{}

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

func (u *AuthUsecase) UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmail) error {
	userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
	if !isExist {
		return errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
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

func (u *AuthUsecase) UpdateUserPassword(ctx context.Context, user dto.UpdateUserPassword) error {
	userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
	if !isExist {
		return errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
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

	passwordHash, err := GeneratePasswordHash(user.NewPassword)
	if err != nil {
		return err
	}

	return u.repo.UpdateUserPassword(ctx, userID, passwordHash)
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
