package repository

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"sync"
)

type UserRepository struct {
	users []models.UserRepo
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: []models.UserRepo{},
	}
}

func (r *UserRepository) CreateUser(user models.UserRepo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = append(r.users, user)

	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.UserRepo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, models.ErrUserNotFound
}

func (r *UserRepository) IncrementUserVersion(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, user := range r.users {
		if user.ID.String() == userID {
			r.users[index].Version++
			return nil
		}
	}

	return models.ErrUserNotFound
}
