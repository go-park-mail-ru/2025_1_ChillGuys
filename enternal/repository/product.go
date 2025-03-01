package repository

import (
	"context"
	"sync"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/enternal/models"
)

type ProductRepo struct{
	storage map[int]*models.Product
	order []int 
	mu sync.RWMutex
}

func NewRepo() *ProductRepo {
	return &ProductRepo{
		storage: make(map[int]*models.Product),
		order: make([]int, 0),
		mu: sync.RWMutex{},
	}
}

func (p *ProductRepo) GetAllPosts(ctx context.Context) ([]*models.Product, error) { //nolint:unparam
	postList := make([]*models.Product, 0, len(p.storage))
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, post := range p.storage {
		postList = append(postList, &(*post))
	}

	return postList, nil
}