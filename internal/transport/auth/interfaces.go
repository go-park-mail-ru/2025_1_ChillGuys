package auth

import "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"

type IUserRepository interface {
	CreateUser(user models.UserRepo) error
	GetUserByEmail(email string) (*models.UserRepo, error)
	IncrementUserVersion(userID string) error
}
