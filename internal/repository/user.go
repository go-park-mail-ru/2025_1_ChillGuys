// map
package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

type UserRepo struct{
	Storage map[int]*models.User
	Mu sync.RWMutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		Storage: make(map[int]*models.User),
		Mu: sync.RWMutex{},
	}
}