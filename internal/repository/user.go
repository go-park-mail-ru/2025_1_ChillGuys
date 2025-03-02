// map
package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

type UserRepo struct{
	storage map[int]*models.User
	mu sync.RWMutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		storage: make(map[int]*models.User),
		mu: sync.RWMutex{},
	}
}