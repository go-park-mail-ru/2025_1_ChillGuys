package repository

import (
	"sync"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

type UserRepository struct {
	users map[string]models.UserDB
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[string]models.UserDB{},
	}
}

func (r *UserRepository) CreateUser(user models.UserDB) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID.String()] = user

	return nil
}

func (r *UserRepository) GetUserCurrentVersion(userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return 0, models.ErrUserNotFound
	}

	return user.Version, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.UserDB, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, models.ErrUserNotFound
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.UserDB, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, ok := r.users[id.String()]; ok {
		return &user, nil
	}

	return nil, models.ErrUserNotFound
}

func (r *UserRepository) IncrementUserVersion(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, isExist := r.users[userID]
	if !isExist {
		return models.ErrUserNotFound
	}

	user.Version++
	r.users[userID] = user

	return nil
}

func (r *UserRepository) CheckUserVersion(userID string, version int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return false
	}

	return user.Version == version
}
