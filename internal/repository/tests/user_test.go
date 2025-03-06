package tests

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Используем repository.NewUserRepository()
	repo := repository.NewUserRepository()
	user := models.UserRepo{
		ID:      uuid.New(), // Генерация UUID с помощью пакета uuid
		Email:   "test@example.com",
		Version: 1,
	}

	// Проверяем, что создание пользователя работает
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Проверяем, что пользователь добавлен
	storedUser, err := repo.GetUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, storedUser.ID)
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Используем repository.NewUserRepository()
	repo := repository.NewUserRepository()
	user := models.UserRepo{
		ID:      uuid.New(), // Генерация UUID с помощью пакета uuid
		Email:   "test@example.com",
		Version: 1,
	}

	// Создаем пользователя
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Проверяем получение пользователя по email
	storedUser, err := repo.GetUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, storedUser.ID)

	// Проверяем несуществующий email
	_, err = repo.GetUserByEmail("nonexistent@example.com")
	assert.Equal(t, models.ErrUserNotFound, err)
}

func TestIncrementUserVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Используем repository.NewUserRepository()
	repo := repository.NewUserRepository()
	user := models.UserRepo{
		ID:      uuid.New(), // Генерация UUID с помощью пакета uuid
		Email:   "test@example.com",
		Version: 1,
	}

	// Создаем пользователя
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Проверяем версию пользователя до инкремента
	version, err := repo.GetUserCurrentVersion(user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 1, version)

	// Инкрементируем версию
	err = repo.IncrementUserVersion(user.ID.String())
	assert.NoError(t, err)

	// Проверяем, что версия инкрементировалась
	version, err = repo.GetUserCurrentVersion(user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 2, version)

	// Проверяем попытку инкремента версии несуществующего пользователя
	err = repo.IncrementUserVersion("nonexistent")
	assert.Equal(t, models.ErrUserNotFound, err)
}

func TestCheckUserVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Используем repository.NewUserRepository()
	repo := repository.NewUserRepository()
	user := models.UserRepo{
		ID:      uuid.New(), // Генерация UUID с помощью пакета uuid
		Email:   "test@example.com",
		Version: 1,
	}

	// Создаем пользователя
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	// Проверяем версию пользователя
	isValid := repo.CheckUserVersion(user.ID.String(), 1)
	assert.True(t, isValid)

	// Проверяем неправильную версию
	isValid = repo.CheckUserVersion(user.ID.String(), 2)
	assert.False(t, isValid)

	// Проверяем несуществующего пользователя
	isValid = repo.CheckUserVersion("nonexistent", 1)
	assert.False(t, isValid)
}
